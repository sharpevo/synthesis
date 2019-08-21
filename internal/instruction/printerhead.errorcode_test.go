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

func TestInstructionPrinterHeadErrorCodeExecute(t *testing.T) {
	ServerNetwork := "tcp"
	ServerAddress := "localhost:21005"
	ricoh_g5.ResetInstance()
	_, err := ricoh_g5.NewDao(ServerNetwork, ServerAddress, 10)
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
			args: []string{ServerAddress},
			expectedRequest: []byte{
				0x01, 0x00, 0x00, 0x00,
			},
			response: []byte{
				0x01, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			args: []string{ServerAddress},
			expectedRequest: []byte{
				0x01, 0x00, 0x00, 0x00,
			},
			response: []byte{
				0x02, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			errString: "response error",
		},
	}
	readyc := make(chan interface{})
	completec := make(chan interface{})

	i := instruction.InstructionPrinterHeadErrorCode{}
	i.Env = interpreter.NewStack()
	deviceVar, _ := vrb.NewVariable(ServerAddress, ServerAddress)
	i.Env.Set(deviceVar)

	go func() {
		<-readyc
		for _, test := range testList {
			_, err := i.Execute(test.args...)
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
			if test.errString == "" && !bytes.Equal(msg, test.expectedRequest) {
				t.Errorf(
					"\nEXPECT: '%x'\nGET:    '%x'\n",
					test.expectedRequest,
					msg,
				)
			}
			t.Logf("Receive mesage: %x\n", msg)
			fmt.Printf("Write mesage: %x\n", test.response)
			conn.Write(test.response)
		}
		conn.Close()
	}()
	<-completec

}

//func TestInstructionPrinterHeadErrorCodeExecuteForRealServer(t *testing.T) {
//t.SkipNow()
//ServerNetwork := "tcp"
//ServerAddress := "192.168.100.215:21005"
//ricoh_g5.ResetInstance()
//_, err := ricoh_g5.NewDao(ServerNetwork, ServerAddress, 1)
//if err != nil {
//t.Fatal(err)
//}

//testList := []struct {
//instruction     instruction.Instructioner
//args            []string
//response        []byte
//expectedRequest []byte
//errString       string
//}{
//{
//instruction: &instruction.InstructionPrinterHeadErrorCode{},
//args:        []string{"var1"},
//},
//{
//instruction: &instruction.InstructionPrinterHeadPrinterStatus{},
//args:        []string{"var1"},
//errString:   "response error",
//},
//}
//for k, test := range testList {
//i := test.instruction
//i.Env = interpreter.NewStack()
//resp, err := i.Execute(test.args...)
//t.Logf("#%d, %v, %v", k, resp, err)
//if err != nil {
//fmt.Println("err", err)
//if test.errString != "" && strings.Contains(err.Error(), test.errString) {
//t.Logf("error occured as expected %s", err)

//continue
//} else {
//// panic if change errString to "foo"
//panic(err)
//}
//}
//v, _ := i.Env.Get(test.args[0])
//actual := v.(*vrb.Variable).Value
//// save to the stack
//if !bytes.Equal(actual.([]byte), resp.([]byte)) {
//t.Errorf(
//"\nEXPECT: '%s'\nGET:    '%x'\n",
//resp,
//actual,
//)
//}
//}
//}
