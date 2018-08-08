package tcp_test

import (
	"fmt"
	"os"
	"posam/protocol/tcp"
	"testing"
)

var p tcp.TCPClient

func TestMain(m *testing.M) {
	p = tcp.TCPClient{
		ServerNetwork: "tcp",
		ServerAddress: "localhost:3333",
	}
	ret := m.Run()
	os.Exit(ret)
}

func TestSend(t *testing.T) {
	for _, i := range []int{4, 5, 6} {
		msg := []byte(fmt.Sprintf("Test-%d", i))
		expected := msg
		actual, err := p.Send(msg, expected)
		if err != nil {
			t.Errorf(
				"\nEXPECT: %q\nGET: %q\n",
				expected,
				actual,
			)

		}
	}
}
