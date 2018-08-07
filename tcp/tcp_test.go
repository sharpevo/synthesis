package tcp_test

import (
	"bytes"
	"fmt"
	"os"
	"posam/protocol/tcp"
	"testing"
)

var p tcp.TCP

func TestMain(m *testing.M) {
	p = tcp.TCP{
		Network: "tcp",
		Address: "localhost:3333",
	}
	ret := m.Run()
	os.Exit(ret)
}

func TestSendString(t *testing.T) {
	for i := range [3]int{} {
		msg := fmt.Sprintf("Test-%d", i)
		actual := p.SendString(msg)
		expected := msg
		if expected != actual {
			t.Errorf(
				"\nEXPECT: %q\nGET: %q\n",
				expected,
				actual,
			)
		}
	}
}

func TestSendByte(t *testing.T) {
	for _, i := range []int{4, 5, 6} {
		msg := []byte(fmt.Sprintf("Test-%d", i))
		actual := p.SendByte(msg)
		expected := msg
		if !bytes.Equal(expected, actual) {
			t.Errorf(
				"\nEXPECT: %q\nGET: %q\n",
				expected,
				actual,
			)

		}
	}
}
