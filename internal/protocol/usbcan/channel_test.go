package usbcan_test

import (
	"fmt"
	"synthesis/protocol/usbcan"
	"reflect"
	"strings"
	"testing"
)

func TestChannelTransmit(t *testing.T) { // {{{
	cases := []struct {
		frameid     int
		recexpected []byte
		recindex    int
		comexpected []byte
		comindex    int

		errmsg      string
		ackResponse *usbcan.Response // pointer can be compared with nil
		comResponse *usbcan.Response // Message is sent as the final resp
	}{
		{ // no ack, no com: e.g., switch, humiture
			frameid:     1,
			recexpected: []byte{},
			recindex:    0,
			comexpected: []byte{},
			comindex:    0,
			errmsg:      "",
			comResponse: &usbcan.Response{
				Message: []byte{1},
			},
		},
		{ // ack com: e.g., switch advanced, motor
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED}, // 0x00
			recindex:    6,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED}, // 0x01
			comindex:    6,
			errmsg:      "",
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_RECEIVED, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_COMPLETED, 9},
			},
		},
		{ // no ack, com: e.g., rom read & write
			frameid:     1,
			recexpected: []byte{},
			recindex:    0,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED},
			comindex:    2,
			errmsg:      "",
			comResponse: &usbcan.Response{
				Message: []byte{0, 0, usbcan.STATUS_CODE_COMPLETED, 0, 0, 0, 0, 0},
			},
		},

		{ // ack com: ack error
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{}, // stop sending extra message to channel
			comindex:    6,
			errmsg:      "ack error",
			ackResponse: &usbcan.Response{
				Error: fmt.Errorf("ack error"),
			},
			comResponse: &usbcan.Response{
				Message: []byte{}, // nil when error occured
			},
		},
		{ // ack com: ack status error
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{}, // stop sending extra message to channel
			comindex:    6,
			errmsg:      "invalid status code 0xe",
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, 0xE, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{}, // nil when error occured
			},
		},
		{ // ack com: com error
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED},
			comindex:    6,
			errmsg:      "com error",
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_RECEIVED, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{}, // nil when error occured
				Error:   fmt.Errorf("com error"),
			},
		},
		{ // ack com: com error status code
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED},
			comindex:    6,
			errmsg:      "unknown error when execute []byte{0x9, 0x9, 0x9, 0x9, 0x9, 0x9, 0x2, 0x9", // omit instruction code
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_RECEIVED, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_ERROR, 9},
			},
		},
		{ // ack com: com invalid status code
			frameid:     1,
			recexpected: []byte{usbcan.STATUS_CODE_RECEIVED},
			recindex:    6,
			comexpected: []byte{usbcan.STATUS_CODE_COMPLETED},
			comindex:    6,
			errmsg:      "invalid status code 0xe", // omit instruction code
			ackResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, usbcan.STATUS_CODE_RECEIVED, 9},
			},
			comResponse: &usbcan.Response{
				Message: []byte{9, 9, 9, 9, 9, 9, 0xE, 9},
			},
		},
	}
	next := make(chan struct{})
	for index, c := range cases {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			channel := &usbcan.Channel{}
			channel.Init()
			go func() {
				resp, err := channel.Transmit(
					c.frameid,
					c.comResponse.Message,
					c.recexpected,
					c.recindex,
					c.comexpected,
					c.comindex,
				)
				if err != nil && !strings.HasPrefix(err.Error(), c.errmsg) {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.errmsg,
						err.Error(),
					)
				}
				if err == nil && !reflect.DeepEqual(resp, c.comResponse.Message) { // not same when error
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.comResponse.Message,
						resp,
					)
				}
				next <- struct{}{}
			}()
			reqi, _ := channel.RequestQueue.Pop()
			req := reqi.(*usbcan.Request)
			if c.ackResponse != nil {
				req.Responsec <- *c.ackResponse
			}
			if c.errmsg == "" || // not sending com when error occured, especially ack error
				len(c.comexpected) > 0 { // still sending when com error expected
				req.Responsec <- *c.comResponse
			}
			<-next
		})
	}
} // }}}
