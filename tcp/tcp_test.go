package tcp_test

import (
	"net"
	"os"
	"posam/protocol/tcp"
	"testing"
	"time"
)

var ServerNetwork = "tcp"
var ServerAddress = "localhost:6507"
var client = tcp.TCPClient{
	ServerNetwork: ServerNetwork,
	ServerAddress: ServerAddress,
}

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}

func TestSend(t *testing.T) {

	testList := []struct {
		timeout     time.Duration
		message     []byte
		expected    []byte
		serverSleep time.Duration
	}{
		{
			timeout:  1 * time.Second,
			message:  []byte("Test-1"),
			expected: []byte("Test-1-expected"),
		},
		{
			timeout:  1 * time.Second,
			message:  []byte("Test-2"),
			expected: []byte("Test-2-messages"),
		},
		{
			timeout:     1 * time.Second,
			message:     []byte("Test-3"),
			expected:    []byte(""),
			serverSleep: 2 * time.Second,
		},
	}
	readyc := make(chan interface{})
	go func() {
		for _, test := range testList {
			<-readyc // send messages after the server is launched
			actual, err := client.Send(test.message, test.expected)
			t.Logf("Send message: %s", string(test.message))
			if err != nil {
				t.Errorf(
					"\nEXPECT: %q\nGET: %q\n",
					test.expected,
					actual,
				)
			}
		}
	}()

	l, err := net.Listen(ServerNetwork, ServerAddress)
	defer l.Close()
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range testList {
		readyc <- true
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, 32)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		msg := buf[:n]
		t.Logf("Receive mesage: %s", msg)
		if test.serverSleep != 0 {
			t.Logf("Server sleep: %s", test.serverSleep)
			time.Sleep(test.serverSleep)
		}
		t.Logf("Write mesage: %s", test.expected)
		conn.Write(test.expected)
		conn.Close()
	}
}
