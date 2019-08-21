package instruction_test

import (
	"bytes"
	"fmt"
	"net"
	"synthesis/dao/ricoh_g5"
	"synthesis/instruction"
	"synthesis/interpreter"
	"synthesis/interpreter/vrb"
	"strings"
	"testing"
)

func TestInstructionPrinterHeadWaveformExecute(t *testing.T) {
	ServerNetwork := "tcp"
	ServerAddress := "localhost:6507"
	ricoh_g5.ResetInstance()
	// TODO: concurrency
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
				ServerAddress,
				"0",
				"1",
				"10.24",
				"1",
				"1.1", "2.2", "2.2", // fall
				"4.4", "5.5", "5.5", // hold
				"7.7", "8.8", "6.6", // rising
				"10.10", "11.11", "11.11", // wait
				"1", // bit
			},
			expectedRequest: []byte{
				0x00, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00,
				0x0a, 0xd7, 0x23, 0x41,
				0x1, 0x0, 0x0, 0x0,

				0xcd, 0xcc, 0x8c, 0x3f,
				0xcd, 0xcc, 0xc, 0x40,
				0xcd, 0xcc, 0xc, 0x40,

				0xcd, 0xcc, 0x8c, 0x40,
				0x0, 0x0, 0xb0, 0x40,
				0x0, 0x0, 0xb0, 0x40,

				0x66, 0x66, 0xf6, 0x40,
				0xcd, 0xcc, 0xc, 0x41,
				0x33, 0x33, 0xd3, 0x40,

				0x9a, 0x99, 0x21, 0x41,
				0x8f, 0xc2, 0x31, 0x41,
				0x8f, 0xc2, 0x31, 0x41,

				0x1, 0x0, 0x0, 0x0,
			},
			response: []byte{
				0x04, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			args: []string{
				ServerAddress,
				"1",
				"2",
				"11.22",
				"3",
				"03", "02", "01", "00",
				"03", "02", "01", "00",
				"03", "02", "01", "00",
				"01",
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

	i := instruction.InstructionPrinterHeadWaveform{}
	i.Env = interpreter.NewStack()
	deviceVar, _ := vrb.NewVariable(ServerAddress, ServerAddress)
	i.Env.Set(deviceVar)

	l, err := net.Listen(ServerNetwork, ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		<-readyc
		for _, test := range testList {
			_, err := i.Execute(test.args...)
			if err != nil {
				if test.errString != "" && strings.Contains(err.Error(), test.errString) {
					fmt.Printf("error occured as expected %s", err)

					// no things to be sent if error occurred
					// send a message to server to unblock l.Accept()
					//ricoh_g5.Instance(ServerAddress).QueryErrorCode()

					l.Close()
					continue
				} else {
					// panic if change errString to "foo"
					panic(err)
				}
			}
		}
		completec <- true
	}()

	req := ricoh_g5.WaveformUnit.Request()

	go func() { // or completec will be blocked by the l.Accept()
		readyc <- true
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}
		for _, test := range testList {

			buf := make([]byte, 1536)
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
