package instruction_test

import (
	"bytes"
	"fmt"
	"net"
	"posam/dao/ricoh_g5"
	"posam/instruction"
	"posam/interpreter"
	"posam/util/concurrentmap"
	"strings"
	"testing"
)

func TestInstructionPrinterHeadPrintDataExecute(t *testing.T) {
	ServerNetwork := "tcp"
	ServerAddress := "localhost:6507"
	ricoh_g5.ResetInstance()
	_, err := ricoh_g5.NewDao(ServerNetwork, ServerAddress, 1)
	if err != nil {
		t.Fatal(err)
	}

	testList := []struct {
		args            []string
		response        []byte
		expectedRequest []byte
		errString       string
	}{
		{
			args: []string{
				"var1",
				"2",
				"32",
				"3",
				"030201",
			},
			expectedRequest: []byte{
				0x02, 0x00, 0x00, 0x00,
				0x20, 0x00, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x00,

				0x03, 0x02, 0x01,
			},
			response: []byte{
				0x03, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			args: []string{
				"var1",
				"1",
				"1280",
				"6",
				"030201",
			},
			expectedRequest: []byte{
				0x03, 0x00, 0x00, 0x00,
			},
			response: []byte{
				0x03, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			errString: "translated with unexpected length",
		},
	}
	readyc := make(chan interface{})
	completec := make(chan interface{})

	i := instruction.InstructionPrinterHeadPrintData{}
	i.Env = concurrentmap.NewConcurrentMap()

	l, err := net.Listen(ServerNetwork, ServerAddress)
	defer l.Close()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		<-readyc
		for _, test := range testList {
			resp, err := i.Execute(test.args...)
			if err != nil {
				if test.errString != "" && strings.Contains(err.Error(), test.errString) {
					t.Logf("error occured as expected %s", err)

					// no things to be sent if error occurred
					// send a message to server to unblock l.Accept()
					l.Close()
					continue
				} else {
					// panic if change errString to "foo"
					panic(err)
				}
			}
			v, _ := i.Env.Get(test.args[0])
			actual := v.(*interpreter.Variable).Value
			// save to the stack
			if !bytes.Equal(actual.([]byte), resp.([]byte)) {
				t.Errorf(
					"\nEXPECT: '%s'\nGET:    '%x'\n",
					resp,
					actual,
				)
			}
			fmt.Printf("%#v\n", v)
		}
		completec <- true
	}()

	req := ricoh_g5.PrintDataUnit.Request()

	go func() {
		readyc <- true
		conn, err := l.Accept()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Fatal(err)
			} else {
				panic(err)
			}
		}
		for _, test := range testList {

			buf := make([]byte, 32)
			n, err := conn.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			msg := buf[:n]
			expected := append(req.Bytes(), test.expectedRequest...)
			if test.errString == "" && !bytes.Equal(msg, expected) {
				t.Errorf(
					"\nEXPECT: '%x'\nGET:    '%x'\n",
					expected,
					msg,
				)
			}
			t.Logf("Receive mesage: %x\n", msg)
			t.Logf("Write mesage: %x\n", test.response)
			conn.Write(test.response)
		}
		conn.Close()
	}()
	<-completec

}
