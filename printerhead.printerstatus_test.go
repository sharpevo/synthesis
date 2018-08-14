package instruction_test

import (
	"bytes"
	"net"
	"posam/dao/ricoh_g5"
	"posam/instruction"
	"posam/interpreter"
	"posam/protocol/tcp"
	"posam/util/concurrentmap"
	"strings"
	"testing"
)

func TestInstructionPrinterHeadPrinterStatusExecute(t *testing.T) {
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
			args: []string{"var1"},
			expectedRequest: []byte{
				0x02, 0x00, 0x00, 0x00,
			},
			response: []byte{
				0x02, 0x00, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x00,
			},
		},
		{
			args: []string{"var1"},
			expectedRequest: []byte{
				0x02, 0x00, 0x00, 0x00,
			},
			response: []byte{
				0x03, 0x00, 0x00, 0x00,
				0x03, 0x00, 0x00, 0x00,
			},
			errString: "response error",
		},
	}
	readyc := make(chan interface{})
	completec := make(chan interface{})

	i := instruction.InstructionPrinterHeadPrinterStatus{}
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
					//ricoh_g5.Instance(ServerAddress).QueryErrorCode()

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
		}
		completec <- true
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
		if test.errString == "" && !bytes.Equal(msg, test.expectedRequest) {
			t.Errorf(
				"\nEXPECT: '%x'\nGET:    '%x'\n",
				test.expectedRequest,
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
