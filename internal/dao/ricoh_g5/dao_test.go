package ricoh_g5_test

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"synthesis/dao"
	"synthesis/dao/ricoh_g5"
	"synthesis/protocol/tcp"
	"testing"
)

type MockTCPClient struct {
	tcp.Client
}

func (m *MockTCPClient) Send(message, expected []byte) (resp []byte, err error) {
	fmt.Println("called")
	return
}

var ServerNetwork = "tcp"
var ServerAddress = "localhost:6507"

func TestQueryFunction(t *testing.T) {

	ricoh_g5.ResetInstance()

	ricoh_g5.NewDao(ServerNetwork, ServerAddress, 5)

	testList := []struct {
		function  func() (interface{}, error)
		expected  dao.CompletedResponse
		response  dao.CompletedResponse
		errString string
	}{
		{
			function: func() (interface{}, error) {
				return ricoh_g5.Instance(ServerAddress).QueryErrorCode()
			},
			expected: ricoh_g5.ErrorCodeUnit.ComResp(),
			response: ricoh_g5.ErrorCodeUnit.ComResp(),
		},
		{
			function: func() (interface{}, error) {
				return ricoh_g5.Instance(ServerAddress).QueryErrorCode()
			},
			expected:  ricoh_g5.ErrorCodeUnit.ComResp(),
			response:  dao.CompletedResponse("test"),
			errString: "response error",
		},
		{
			function: func() (interface{}, error) {
				return ricoh_g5.Instance(ServerAddress).QueryPrinterStatus()
			},
			expected: ricoh_g5.PrinterStatusUnit.ComResp(),
			response: ricoh_g5.PrinterStatusUnit.ComResp(),
		},
		{
			function: func() (interface{}, error) {
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
		<-readyc
		for i, test := range testList {
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

			if !bytes.Equal(test.expected, actual.([]byte)) {
				t.Errorf(
					"\nEXPECT: '%x'\nGET: '%x'\n",
					test.expected,
					actual,
				)
			}
		}
		completec <- true
	}()

	l, err := net.Listen(ServerNetwork, ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	readyc <- true
	conn, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}

	for k, test := range testList {
		//readyc <- true
		fmt.Println("#", k)
		buf := make([]byte, 32)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		msg := buf[:n]
		t.Logf("Receive mesage: %x", msg)
		t.Logf("Write mesage: %x", test.response)
		conn.Write(test.response)
	}
	conn.Close()

	<-completec // allow failure in goroutine then complete the test case
	l.Close()
}

func TestPrintData(t *testing.T) {

	ricoh_g5.ResetInstance()
	ricoh_g5.NewDao(ServerNetwork, ServerAddress, 5)

	testList := []struct {
		bitsPerPixel    string
		width           string
		lineBufferSize  string
		lineBuffer      string
		expectedRequest []byte
		expected        []byte
		response        []byte
		errString       string
		ignoreResponse  bool
	}{
		//{
		//bitsPerPixel:   "2",
		//width:          "32",
		//lineBufferSize: "3",
		//lineBuffer:     "030201",
		//expectedRequest: []byte{
		//0x02, 0x00, 0x00, 0x00,
		//0x20, 0x00, 0x00, 0x00,
		//0x03, 0x00, 0x00, 0x00,

		//0x03, 0x02, 0x01,
		//},
		//expected: []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		//response: ricoh_g5.PrintDataUnit.ComResp(),
		//},
		{
			bitsPerPixel:   "1",    // 1 or 2
			width:          "1280", // 320 * 4 = 1280
			lineBufferSize: "160",  // 1280 bits / 8 = 160 bytes
			lineBuffer:     "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expectedRequest: []byte{
				0x01, 0x00, 0x00, 0x00,
				0x00, 0x05, 0x00, 0x00,
				0xa0, 0x00, 0x00, 0x00,

				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			expected: []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			response: ricoh_g5.PrintDataUnit.ComResp(),
		},
		{
			bitsPerPixel:   "1",    // 1 or 2
			width:          "1280", // 320 * 4 = 1280
			lineBufferSize: "160",  // 1280 bits / 8 = 160 bytes
			lineBuffer:     "01020304",
			expectedRequest: []byte{
				0x01, 0x00, 0x00, 0x00,
				0x00, 0x05, 0x00, 0x00,
				0xa0, 0x00, 0x00, 0x00,

				0x01, 0x02, 0x03, 0x04,
			},
			expected:       []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			response:       []byte("failed"),
			errString:      "is translated with unexpected length 4 (160)",
			ignoreResponse: true,
		},
	}

	readyc := make(chan interface{})
	completec := make(chan interface{})

	go func() {
		<-readyc
		for i, test := range testList {
			fmt.Printf(">>%d\n", i)
			actual, err := ricoh_g5.Instance(ServerAddress).PrintData(
				test.bitsPerPixel,
				test.width,
				test.lineBufferSize,
				test.lineBuffer,
				test.lineBuffer,
			)
			if err != nil {
				if test.errString != "" && strings.Contains(err.Error(), test.errString) {
					t.Logf("error occured as expected %s", err)
					//ricoh_g5.Instance(ServerAddress).QueryErrorCode()
					continue
				} else {
					// panic if change errString to "foo"
					panic(err)
				}
			}

			if !bytes.Equal(test.expected, actual.([]byte)) {
				t.Errorf(
					"\nEXPECT: '%x'\nGET: '%x'\n",
					test.expected,
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
	req := ricoh_g5.PrintDataUnit.Request()

	readyc <- true
	conn, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range testList {
		if test.ignoreResponse {
			fmt.Println("response ignored")
			continue
		}
		buf := make([]byte, 256)
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
		t.Logf("Receive mesage: %x", msg)
		t.Logf("Write mesage: %x", test.response)
		conn.Write(test.response)
	}
	conn.Close()

	<-completec // allow failure in goroutine then complete the test case
}

func TestSendWaveform(t *testing.T) {

	ricoh_g5.ResetInstance()
	ricoh_g5.NewDao(ServerNetwork, ServerAddress, 5)

	testList := []struct {
		headBoardIndex      string
		rowIndexOfHeadBoard string
		voltagePercentage   string
		segmentCount        string
		segment             []string
		expected            []byte
		expectedRequest     []byte
		expectedResponse    []byte
		errString           string
		ignoreResponse      bool
	}{
		{
			headBoardIndex:      "0", // 0 for the first head baord
			rowIndexOfHeadBoard: "1", // 0 for the first row of head board
			voltagePercentage:   "10.24",
			segmentCount:        "1",
			segment: []string{
				"1.1", "2.2", "2.2", // fall
				"4.4", "5.5", "5.5", // hold
				"7.7", "8.8", "6.6", // rising
				"10.10", "11.11", "11.11", // wait
				"1", // bit
			},
			expectedRequest: []byte{
				0x0, 0x0, 0x0, 0x0,
				0x1, 0x0, 0x0, 0x0,
				0xa, 0xd7, 0x23, 0x41,
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
			expectedResponse: ricoh_g5.WaveformUnit.ComResp(),
		},
		{
			headBoardIndex:      "1", // 0 for the first head baord
			rowIndexOfHeadBoard: "2", // 0 for the first row of head board
			voltagePercentage:   "11.22",
			segmentCount:        "3",
			segment: []string{
				"1.1", "2.2", "3.3", // fall
				"4.4", "5.5", "6.6", // hold
				"7.7", "8.8", "9.9", // rising
				"10.10", "11.11", "12.12", // wait
				"1", // bit
			},
			expectedRequest: []byte{
				0x01, 0x00, 0x00, 0x00,
				0x02, 0x00, 0x00, 0x00,
				0x1f, 0x85, 0x33, 0x41,
				0x03, 0x00, 0x00, 0x00,

				0xcd, 0xcc, 0x8c, 0x3f,
				0xcd, 0xcc, 0xc, 0x40,
				0x33, 0x33, 0x53, 0x40,

				0xcd, 0xcc, 0x8c, 0x40,
				0x0, 0x0, 0xb0, 0x40,
				0x33, 0x33, 0xd3, 0x40,

				0x66, 0x66, 0xf6, 0x40,
				0xcd, 0xcc, 0xc, 0x41,
				0x66, 0x66, 0x1e, 0x41,

				0x9a, 0x99, 0x21, 0x41,
				0x8f, 0xc2, 0x31, 0x41,
				0x85, 0xeb, 0x41, 0x41,

				0x1, 0x0, 0x0, 0x0,
			},
			expectedResponse: ricoh_g5.WaveformUnit.ComResp(),
			errString:        "translated with unexpected length",
			ignoreResponse:   true,
		},
	}

	readyc := make(chan interface{})
	completec := make(chan interface{})

	go func() {
		<-readyc
		for i, test := range testList {
			//<-readyc
			fmt.Printf(">>%d\n", i)
			actual, err := ricoh_g5.Instance(ServerAddress).SendWaveform(
				test.headBoardIndex,
				test.rowIndexOfHeadBoard,
				test.voltagePercentage,
				test.segmentCount,
				test.segment,
			)
			if err != nil {
				if test.errString != "" && strings.Contains(err.Error(), test.errString) {
					fmt.Printf("error occured as expected %s", err)

					// no things to be sent if error occurred
					// send a message to server to unblock l.Accept()
					//ricoh_g5.Instance(ServerAddress).QueryErrorCode()

					continue
				} else {
					// panic if change errString to "foo"
					panic(err)
				}
			}

			if !bytes.Equal(test.expectedResponse, actual.([]byte)) {
				t.Errorf(
					"\nEXPECT: '%x'\nGET: '%x'\n",
					test.expected,
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

	req := ricoh_g5.WaveformUnit.Request()

	readyc <- true
	conn, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range testList {
		if test.ignoreResponse {
			fmt.Println("response ignored")
			continue
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		msg := buf[:n]

		expected := append(req.Bytes(), test.expectedRequest...)
		if test.errString == "" && !bytes.Equal(msg, expected) {
			t.Errorf(
				"\nEXPECT: '%#v'\nGET:    '%#v'\n",
				expected,
				msg,
			)
		}
		t.Logf("Receive mesage: %x", msg)
		t.Logf("Write mesage: %x", test.expectedResponse)
		conn.Write(test.expectedResponse)
	}
	conn.Close()

	<-completec // allow failure in goroutine then complete the test case
}

func TestSegmentifyAndSegmentBytes(t *testing.T) {
	testList := []struct {
		argumentList []string
		expected     []byte
		errString    string
	}{
		{
			argumentList: []string{
				"1.1", "2.2", "3.3", // fall
				"4.4", "3.3", "3.3", // hold
				"7.7", "3.3", "5.5", // rising
				"10.10", "5.5", "5.5", // wait
				"1", // bit

				//"2.1", "3.2", "4.3", // fall
				//"2.4", "3.5", "4.6", // hold
				//"2.7", "3.8", "4.9", // rising
				//"20.10", "31.11", "42.12", // wait
				//"2", // bit
			},
			expected: []byte{
				0xcd, 0xcc, 0x8c, 0x3f, // 1.1
				0xcd, 0xcc, 0xc, 0x40, // 2.2
				0x33, 0x33, 0x53, 0x40, // 3.3

				0xcd, 0xcc, 0x8c, 0x40,
				0x33, 0x33, 0x53, 0x40,
				0x33, 0x33, 0x53, 0x40,

				0x66, 0x66, 0xf6, 0x40,
				0x33, 0x33, 0x53, 0x40,
				0x0, 0x0, 0xb0, 0x40,

				0x9a, 0x99, 0x21, 0x41,
				0x0, 0x0, 0xb0, 0x40,
				0x0, 0x0, 0xb0, 0x40,

				0x1, 0x0, 0x0, 0x0,

				//0xcd, 0xcc, 0x8c, 0x3f,
				//0xcd, 0xcc, 0xc, 0x40,
				//0x33, 0x33, 0x53, 0x40,

				//0xcd, 0xcc, 0x8c, 0x40,
				//0x0, 0x0, 0xb0, 0x40,
				//0x33, 0x33, 0xd3, 0x40,

				//0x66, 0x66, 0xf6, 0x40,
				//0xcd, 0xcc, 0xc, 0x41,
				//0x66, 0x66, 0x1e, 0x41,

				//0x9a, 0x99, 0x21, 0x41,
				//0x8f, 0xc2, 0x31, 0x41,
				//0x85, 0xeb, 0x41, 0x41,

				//0x1, 0x0, 0x0, 0x0,

				//0x66, 0x66, 0x6, 0x40,
				//0xcd, 0xcc, 0x4c, 0x40,
				//0x9a, 0x99, 0x89, 0x40,

				//0x9a, 0x99, 0x19, 0x40,
				//0x0, 0x0, 0x60, 0x40,
				//0x33, 0x33, 0x93, 0x40,

				//0xcd, 0xcc, 0x2c, 0x40,
				//0x33, 0x33, 0x73, 0x40,
				//0xcd, 0xcc, 0x9c, 0x40,

				//0xcd, 0xcc, 0xa0, 0x41,
				//0x48, 0xe1, 0xf8, 0x41,
				//0xe1, 0x7a, 0x28, 0x42,

				//0x2, 0x0, 0x0, 0x0,
			},
		},
		{
			argumentList: []string{
				"1.1", "2.2", "3.3", // fall
				"4.4", "5.5", "6.6", // hold
				"7.7", "8.8", "9.9", // rising
				"10.10", "11.11", "12.12", // wait
				"1", // bit

				"2.1", "3.2", "4.3", // fall
				"2.4", "3.5", "4.6", // hold
				"2.7", "3.8", "4.9", // rising
				"20.10", "31.11", "42.12", // wait
			},
			expected:  []byte{},
			errString: "invalid segment",
		},
	}

	for _, test := range testList {
		segmentList, err := ricoh_g5.Segmentify(test.argumentList, 13)
		if err != nil {
			if strings.Contains(err.Error(), test.errString) {
				fmt.Println("error occurred as expected", err)
			} else {
				t.Errorf(err.Error())
			}
		}
		segmentBytes, err := ricoh_g5.SegmentBytes(segmentList, 2)
		if err != nil {
			t.Errorf(err.Error())
		}
		if !bytes.Equal(segmentBytes, test.expected) {
			t.Errorf(
				"\nEXPECT: '%#v'\nGET:    '%#v'\n",
				test.expected,
				segmentBytes,
			)
		}
	}
}
