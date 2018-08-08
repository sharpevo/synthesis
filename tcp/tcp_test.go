package tcp_test

import (
	"fmt"
	"net"
	"os"
	"posam/protocol/tcp"
	"testing"
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
	for _, i := range []int{4, 5, 6} {

		message := []byte(fmt.Sprintf("Test-%d", i))
		expected := append(message, []byte("-processed")...)

		go func() {

			actual, err := client.Send(message, expected)
			t.Logf("Send message: %s", string(message))
			if err != nil {
				t.Errorf(
					"\nEXPECT: %q\nGET: %q\n",
					expected,
					actual,
				)
			}
		}()

	}

	l, err := net.Listen(ServerNetwork, ServerAddress)
	defer l.Close()
	if err != nil {
		t.Fatal(err)
	}
	for _ = range [3]int{} {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		msg := string(buf[:n])
		t.Logf("Receive mesage: %s", msg)
		resp := append(buf[:n], []byte("-processed")...)
		conn.Write(resp)
		conn.Close()
	}
}
