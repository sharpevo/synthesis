package tcp_test

import (
	"net"
	"os"
	"posam/protocol/tcp"
	"testing"
	"time"
)

type MockConnectioner struct {
	tcp.Connectivity
}

func (m *MockConnectioner) Connect(network string, address string, timeout time.Duration) (conn *net.TCPConn, err error) {
	return
}

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}

func TestSendSerial(t *testing.T) {

	ServerNetwork := "tcp"
	ServerAddress := "localhost:6507"
	client := tcp.TCPClient{
		Connectivitier:    &tcp.Connectivity{},
		ServerNetwork:     ServerNetwork,
		ServerAddress:     ServerAddress,
		ServerTimeout:     1 * time.Second,
		ServerConcurrency: false,
	}

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
	completec := make(chan interface{})
	go func() {
		<-readyc // send messages after the server is launched
		for i, test := range testList {
			actual, err := client.Send(test.message, test.expected)
			t.Logf("#%d Send message: %s", i, string(test.message))
			if err != nil {
				if test.serverSleep == 0 {
					t.Errorf(
						"\nEXPECT: %q\nGET: %q\n",
						test.expected,
						actual,
					)
				} else {
					t.Logf("timeout error occured as expected: %v", err)
				}
			}
		}
		completec <- true
	}()

	l, err := net.Listen(ServerNetwork, ServerAddress)
	defer l.Close()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		readyc <- true

		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}
		for _, test := range testList {
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
		}

		conn.Close()
	}()
	<-completec
}

func TestSendConcurrent(t *testing.T) {

	ServerNetwork := "tcp"
	ServerAddress := "localhost:6508"
	client := tcp.TCPClient{
		Connectivitier:    &tcp.Connectivity{},
		ServerNetwork:     ServerNetwork,
		ServerAddress:     ServerAddress,
		ServerTimeout:     1 * time.Second,
		ServerConcurrency: true,
	}

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
	completec := make(chan interface{})
	go func() {
		<-readyc // send messages after the server is launched
		for i, test := range testList {
			actual, err := client.Send(test.message, test.expected)
			t.Logf("#%d Send message: %s", i, string(test.message))
			if err != nil {
				if test.serverSleep == 0 {
					t.Errorf(
						"\nEXPECT: %q\nGET: %q\n",
						test.expected,
						actual,
					)
				} else {
					t.Logf("timeout error occured as expected: %v", err)
				}
			}
		}
		completec <- true
	}()

	l, err := net.Listen(ServerNetwork, ServerAddress)
	defer l.Close()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		readyc <- true

		for _, test := range testList {
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
	}()
	<-completec
}
