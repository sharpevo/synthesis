package ricoh_g5_test

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"posam/dao"
	"posam/dao/ricoh_g5"
	"posam/protocol/tcp"
	"strings"
	"testing"
)

type MockTCPClient struct {
	tcp.TCPClient
}

func (m *MockTCPClient) Send(message, expected []byte) (resp []byte, err error) {
	fmt.Println("called")
	return
}

var ServerNetwork = "tcp"
var ServerAddress = "localhost:6507"

func TestMain(m *testing.M) {

	ricohDao := &ricoh_g5.Dao{
		DeviceAddress: ServerAddress,
		TCPClient: &tcp.TCPClient{
			Connectioner:  &tcp.Connection{},
			ServerNetwork: ServerNetwork,
			ServerAddress: ServerAddress,
		},
	}
	ricoh_g5.AddInstance(ricohDao)
	ret := m.Run()
	os.Exit(ret)
}

func TestQueryFunction(t *testing.T) {
	testList := []struct {
		function  func() (string, error)
		expected  dao.CompletedResponse
		response  dao.CompletedResponse
		errString string
	}{
		{
			function: func() (string, error) {
				return ricoh_g5.Instance(ServerAddress).QueryErrorCode()
			},
			expected: ricoh_g5.ErrorCodeUnit.ComResp(),
			response: ricoh_g5.ErrorCodeUnit.ComResp(),
		},
		{
			function: func() (string, error) {
				return ricoh_g5.Instance(ServerAddress).QueryErrorCode()
			},
			expected:  ricoh_g5.ErrorCodeUnit.ComResp(),
			response:  dao.CompletedResponse("test"),
			errString: "response error",
		},
		{
			function: func() (string, error) {
				return ricoh_g5.Instance(ServerAddress).QueryPrinterStatus()
			},
			expected: ricoh_g5.PrinterStatusUnit.ComResp(),
			response: ricoh_g5.PrinterStatusUnit.ComResp(),
		},
		{
			function: func() (string, error) {
				return ricoh_g5.Instance(ServerAddress).QueryPrinterStatus()
			},
			expected:  ricoh_g5.PrinterStatusUnit.ComResp(),
			response:  dao.CompletedResponse("test"),
			errString: "response error",
		},
	}

	readyc := make(chan interface{})
	completec := make(chan interface{})
	go func() {
		for i, test := range testList {
			<-readyc
			t.Logf(">>%d", i)
			actual, err := test.function()
			if err != nil {
				if test.errString != "" && strings.Contains(err.Error(), test.errString) {
					t.Logf("error occured as expected %s", err)
					continue // omit the expectation checking
				} else {
					// panic if change errString to "foo"
					panic(err)
				}
			}

			if !bytes.Equal(test.expected, []byte(actual)) {
				t.Errorf(
					"\nEXPECT: '%x'\nGET: '%x'\n",
					test.expected,
					[]byte(actual),
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
		t.Logf("Receive mesage: %x", msg)
		t.Logf("Write mesage: %x", test.response)
		conn.Write(test.response)
		conn.Close()
	}

	<-completec // allow failure in goroutine then complete the test case
}

func TestPrintData(t *testing.T) {
	testList := []struct {
		bitsPerPixel   []byte
		width          []byte
		lineBufferSize []byte
		lineBuffer     []byte
		expected       []byte
		response       []byte
		errString      string
	}{
		{
			bitsPerPixel:   []byte{0x01, 0x00, 0x00, 0x00}, // 0x01 or 0x02
			width:          []byte{0x00, 0x05, 0x00, 0x00}, // 320 * 4 = 1280
			lineBufferSize: []byte{0xA0, 0x00, 0x00, 0x00}, // 1280 bits/160 bytes
			lineBuffer:     []byte{},
			expected:       []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			response:       ricoh_g5.PrintDataUnit.ComResp(),
		},
		{
			bitsPerPixel:   []byte{0x01, 0x00, 0x00, 0x00}, // 0x01 or 0x02
			width:          []byte{0x00, 0x05, 0x00, 0x00}, // 320 * 4 = 1280
			lineBufferSize: []byte{0xA0, 0x00, 0x00, 0x00}, // 1280 bits/160 bytes
			lineBuffer:     []byte{},
			expected:       []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			response:       []byte("failed"),
			errString:      "response error",
		},
	}

	readyc := make(chan interface{})
	completec := make(chan interface{})

	go func() {
		for i, test := range testList {
			<-readyc
			t.Logf(">>%d", i)
			actual, err := ricoh_g5.Instance(ServerAddress).PrintData(
				test.bitsPerPixel,
				test.width,
				test.lineBufferSize,
				test.lineBuffer,
			)
			if err != nil {
				if test.errString != "" && strings.Contains(err.Error(), test.errString) {
					t.Logf("error occured as expected %s", err)
					continue
				} else {
					// panic if change errString to "foo"
					panic(err)
				}
			}

			if !bytes.Equal(test.expected, []byte(actual)) {
				t.Errorf(
					"\nEXPECT: '%x'\nGET: '%x'\n",
					test.expected,
					[]byte(actual),
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
		t.Logf("Receive mesage: %x", msg)
		t.Logf("Write mesage: %x", test.response)
		conn.Write(test.response)
		conn.Close()
	}

	<-completec // allow failure in goroutine then complete the test case
}
