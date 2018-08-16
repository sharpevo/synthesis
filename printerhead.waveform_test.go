package instruction_test

import (
	"bytes"
	"fmt"
	"net"
	"posam/dao/ricoh_g5"
	"posam/instruction"
	"posam/interpreter"
	"posam/protocol/tcp"
	"posam/util/concurrentmap"
	"strings"
	"testing"
)

func TestInstructionPrinterHeadWaveformExecute(t *testing.T) {
	ServerNetwork := "tcp"
	ServerAddress := "localhost:6507"
	ricoh_g5.ResetInstance()
	ricoh_g5.AddInstance(&ricoh_g5.Dao{
		DeviceAddress: ServerAddress,
		TCPClient: &tcp.TCPClient{
			Connectivitier:    &tcp.Connectivity{},
			ServerNetwork:     ServerNetwork,
			ServerAddress:     ServerAddress,
			ServerConcurrency: true,
		},
	})
	testList := []struct {
		args            []string
		response        []byte
		expectedRequest []byte
		errString       string
	}{
		{
			args: []string{
				"var1",
				"0",
				"1",
				"10.24",
				"5",
				"0302010000",
			},
			expectedRequest: []byte{
				0x00, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00,
				0x0a, 0xd7, 0x23, 0x41,
				0x05, 0x00, 0x00, 0x00,
				0x03, 0x02, 0x01, 0x00, 0x00,
			},
			response: []byte{
				0x04, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			args: []string{
				"var1",
				"1",
				"2",
				"11.22",
				"3",
				"0302010000",
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
	i.Env = concurrentmap.NewConcurrentMap()

	go func() {
		for _, test := range testList {
			<-readyc
			resp, err := i.Execute(test.args...)
			if err != nil {
				if test.errString != "" && strings.Contains(err.Error(), test.errString) {
					t.Logf("error occured as expected %s", err)

					// no things to be sent if error occurred
					// send a message to server to unblock l.Accept()
					ricoh_g5.Instance(ServerAddress).QueryErrorCode()

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

	l, err := net.Listen(ServerNetwork, ServerAddress)
	defer l.Close()
	if err != nil {
		t.Fatal(err)
	}
	req := ricoh_g5.WaveformUnit.Request()
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
		conn.Close()
	}
	<-completec

}
