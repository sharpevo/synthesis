package canalystii

import (
	"fmt"
	"posam/util/concurrentmap"
	"reflect"
	"strings"
	"testing"
)

func TestAddInstance(t *testing.T) { // {{{
	deviceMap = concurrentmap.NewConcurrentMap()
	d := &Dao{_id: "id"}
	fmt.Println(d)
	addInstance(d)
	if _, found := deviceMap.Get("id"); !found {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			"instance found",
			"not found",
		)
	}
} // }}}

func TestSetID(t *testing.T) { // {{{
	cases := []struct {
		id     string
		errmsg string
	}{
		{
			"id",
			"",
		},
		{
			"id",
			"is duplicated",
		},
	}
	deviceMap = concurrentmap.NewConcurrentMap()
	d := &Dao{}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err := d.setID(c.id)
			if c.errmsg != "" {
				if err == nil {
					t.Errorf("expect error: %v\n", c.errmsg)
					return
				}
				if !strings.Contains(err.Error(), c.errmsg) {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.errmsg,
						err.Error(),
					)
				}
			}
			expect := c.id
			actual := d.id()
			if actual != expect {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					expect,
					actual,
				)
			}
		})
	}
} // }}}

func TestMoveRelative(t *testing.T) { // {{{
	cases := []struct {
		motorcode int
		direction int
		speed     int
		position  int

		message []byte
		recresp []byte
		comresp []byte
		output  []byte
		resp    uint16
		err     error
	}{
		{
			1, 2, 3, 4,
			[]byte{
				MotorMoveRelativeUnit.Request().Function,
				1, 2, 0, 3, 0, 4,
			},
			MotorMoveRelativeUnit.RecResp(),
			MotorMoveRelativeUnit.ComResp(),
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			1029, // binary.BigEndian.Uint16([]byte{4, 5}), 00000100,00000101
			nil,
		},
		{
			256, 2, 3, 4,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			0,
			fmt.Errorf("256 overflows uint8"),
		},
		{
			1, 256, 3, 4,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			0,
			fmt.Errorf("256 overflows uint8"),
		},
		{
			1, 2, 65536, 4,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			0,
			fmt.Errorf("65536 overflows uint16"),
		},
		{
			1, 2, 3, 65536,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			0,
			fmt.Errorf("65536 overflows uint16"),
		},
		{
			1, 2, 3, 4,
			[]byte{
				MotorMoveRelativeUnit.Request().Function,
				1, 2, 0, 3, 0, 4,
			},
			MotorMoveRelativeUnit.RecResp(),
			MotorMoveRelativeUnit.ComResp(),
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			1029, // binary.BigEndian.Uint16([]byte{4, 5}), 00000100,00000101
			fmt.Errorf("some error"),
		},
	}
	originSendAck2 := sendAck2
	defer func() { sendAck2 = originSendAck2 }()
	for i, c := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			d := &Dao{}
			sendAck2 = func(
				d *Dao,
				message []byte,
				recResp []byte,
				comResp []byte,
			) ([]byte, error) {
				if !reflect.DeepEqual(message, c.message) ||
					!reflect.DeepEqual(recResp, c.recresp) ||
					!reflect.DeepEqual(comResp, c.comresp) {
					t.Errorf(
						"\nEXPECT: %v %v %v\n GET: %v %v %v\n\n",
						c.message, c.recresp, c.comresp,
						message, recResp, comResp,
					)
				}
				return c.output, c.err
			}
			resp, err := d.MoveRelative(
				c.motorcode,
				c.direction,
				c.speed,
				c.position,
			)
			if err != nil && c.err == nil {
				t.Fatal(err)
			}
			if err != nil && !strings.Contains(err.Error(), c.err.Error()) {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.err.Error(),
					err.Error(),
				)
			}
			if err == nil {
				actual := resp.(uint16)
				if actual != c.resp {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.resp,
						actual,
					)
				}
			}
		})
	}
} // }}}

func TestMoveAbsolute(t *testing.T) { // {{{
	cases := []struct {
		motorcode int
		position  int

		message []byte
		recresp []byte
		comresp []byte
		output  []byte
		resp    uint16
		err     error
	}{
		{
			1, 2,
			[]byte{
				MotorMoveAbsoluteUnit.Request().Function,
				1, 0, 2, 0, 0, 0,
			},
			MotorMoveAbsoluteUnit.RecResp(),
			MotorMoveAbsoluteUnit.ComResp(),
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			1029, // binary.BigEndian.Uint16([]byte{4, 5}), 00000100,00000101
			nil,
		},
		{
			256, 2,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			0,
			fmt.Errorf("256 overflows uint8"),
		},
		{
			1, 65536,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			0,
			fmt.Errorf("65536 overflows uint16"),
		},
		{
			1, 2,
			[]byte{
				MotorMoveAbsoluteUnit.Request().Function,
				1, 0, 2, 0, 0, 0,
			},
			MotorMoveAbsoluteUnit.RecResp(),
			MotorMoveAbsoluteUnit.ComResp(),
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			1029, // binary.BigEndian.Uint16([]byte{4, 5}), 00000100,00000101
			fmt.Errorf("some error"),
		},
	}
	originSendAck2 := sendAck2
	defer func() { sendAck2 = originSendAck2 }()
	for i, c := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			d := &Dao{}
			sendAck2 = func(
				d *Dao,
				message []byte,
				recResp []byte,
				comResp []byte,
			) ([]byte, error) {
				if !reflect.DeepEqual(message, c.message) ||
					!reflect.DeepEqual(recResp, c.recresp) ||
					!reflect.DeepEqual(comResp, c.comresp) {
					t.Errorf(
						"\nEXPECT: %v %v %v\n GET: %v %v %v\n\n",
						c.message, c.recresp, c.comresp,
						message, recResp, comResp,
					)
				}
				return c.output, c.err
			}
			resp, err := d.MoveAbsolute(
				c.motorcode,
				c.position,
			)
			if err != nil && c.err == nil {
				t.Fatal(err)
			}
			if err != nil && !strings.Contains(err.Error(), c.err.Error()) {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.err.Error(),
					err.Error(),
				)
			}
			if err == nil {
				actual := resp.(uint16)
				if actual != c.resp {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.resp,
						actual,
					)
				}
			}
		})
	}
} // }}}

func TestResetMotor(t *testing.T) { // {{{
	cases := []struct {
		motorcode int
		direction int

		message []byte
		recresp []byte
		comresp []byte
		output  []byte
		resp    []byte
		err     error
	}{
		{
			1, 2,
			[]byte{
				MotorResetUnit.Request().Function,
				1, 2, 0, 0, 0, 0,
			},
			MotorResetUnit.RecResp(),
			MotorResetUnit.ComResp(),
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			nil,
		},
		{
			256, 2,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{1},
			fmt.Errorf("256 overflows uint8"),
		},
		{
			1, 256,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{1},
			fmt.Errorf("256 overflows uint8"),
		},
		{
			1, 2,
			[]byte{
				MotorResetUnit.Request().Function,
				1, 2, 0, 0, 0, 0,
			},
			MotorResetUnit.RecResp(),
			MotorResetUnit.ComResp(),
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			fmt.Errorf("some error"),
		},
	}
	originSendAck2 := sendAck2
	defer func() { sendAck2 = originSendAck2 }()
	for i, c := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			d := &Dao{}
			sendAck2 = func(
				d *Dao,
				message []byte,
				recResp []byte,
				comResp []byte,
			) ([]byte, error) {
				if !reflect.DeepEqual(message, c.message) ||
					!reflect.DeepEqual(recResp, c.recresp) ||
					!reflect.DeepEqual(comResp, c.comresp) {
					t.Errorf(
						"\nEXPECT: %v %v %v\n GET: %v %v %v\n\n",
						c.message, c.recresp, c.comresp,
						message, recResp, comResp,
					)
				}
				return c.output, c.err
			}
			resp, err := d.ResetMotor(
				c.motorcode,
				c.direction,
			)
			if err != nil && c.err == nil {
				t.Fatal(err)
			}
			if err != nil && !strings.Contains(err.Error(), c.err.Error()) {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.err.Error(),
					err.Error(),
				)
			}
			if err == nil {
				if !reflect.DeepEqual(resp, c.resp) {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.resp,
						resp,
					)
				}
			}
		})
	}
} // }}}

func TestControlSwitcher(t *testing.T) { // {{{
	cases := []struct {
		data int

		message []byte
		output  []byte
		resp    []byte
		err     error
	}{
		{
			1,
			[]byte{
				SwitcherControlUnit.Request().Function,
				0, 1, 0, 0, 0, 0,
			},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			nil,
		},
		{
			65536,
			[]byte{},
			[]byte{},
			[]byte{},
			fmt.Errorf("65536 overflows uint16"),
		},
		{
			1,
			[]byte{
				SwitcherControlUnit.Request().Function,
				0, 1, 0, 0, 0, 0,
			},
			[]byte{},
			[]byte{},
			fmt.Errorf("some error"),
		},
	}
	originSend := send
	defer func() { send = originSend }()
	for i, c := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			d := &Dao{}
			send = func(
				d *Dao,
				message []byte,
			) ([]byte, error) {
				if !reflect.DeepEqual(message, c.message) {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.message,
						message,
					)
				}
				return c.output, c.err
			}
			resp, err := d.ControlSwitcher(c.data)
			if err != nil && c.err == nil {
				t.Fatal(err)
			}
			if err != nil && !strings.Contains(err.Error(), c.err.Error()) {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.err.Error(),
					err.Error(),
				)
			}
			if err == nil {
				if !reflect.DeepEqual(resp, c.resp) {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.resp,
						resp,
					)
				}
			}
		})
	}
} // }}}

func TestControlSwitcherAdvanced(t *testing.T) { // {{{
	cases := []struct {
		data  int
		speed int
		count int

		message []byte
		recresp []byte
		comresp []byte
		output  []byte
		resp    []byte
		err     error
	}{
		{
			1, 2, 3,
			[]byte{
				SwitcherControlAdvancedUnit.Request().Function,
				0, 1, 2, 0, 3, 0,
			},
			SwitcherControlAdvancedUnit.RecResp(),
			SwitcherControlAdvancedUnit.ComResp(),
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			nil,
		},
		{
			65536, 2, 3,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			fmt.Errorf("65536 overflows uint16"),
		},
		{
			1, 256, 3,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			fmt.Errorf("256 overflows uint8"),
		},
		{
			1, 2, 65536,
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			[]byte{},
			fmt.Errorf("65536 overflows uint16"),
		},
		{
			1, 2, 3,
			[]byte{
				SwitcherControlAdvancedUnit.Request().Function,
				0, 1, 2, 0, 3, 0,
			},
			SwitcherControlAdvancedUnit.RecResp(),
			SwitcherControlAdvancedUnit.ComResp(),
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			fmt.Errorf("some error"),
		},
	}
	originSendAck6 := sendAck6
	defer func() { sendAck6 = originSendAck6 }()
	for i, c := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			d := &Dao{}
			sendAck6 = func(
				d *Dao,
				message []byte,
				recResp []byte,
				comResp []byte,
			) ([]byte, error) {
				if !reflect.DeepEqual(message, c.message) ||
					!reflect.DeepEqual(recResp, c.recresp) ||
					!reflect.DeepEqual(comResp, c.comresp) {
					t.Errorf(
						"\nEXPECT: %v %v %v\n GET: %v %v %v\n\n",
						c.message, c.recresp, c.comresp,
						message, recResp, comResp,
					)
				}
				return c.output, c.err
			}
			resp, err := d.ControlSwitcherAdvanced(
				c.data,
				c.speed,
				c.count,
			)
			if err != nil && c.err == nil {
				t.Fatal(err)
			}
			if err != nil && !strings.Contains(err.Error(), c.err.Error()) {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.err.Error(),
					err.Error(),
				)
			}
			if err == nil {
				if !reflect.DeepEqual(resp, c.resp) {
					t.Errorf(
						"\nEXPECT: %v\n GET: %v\n\n",
						c.resp,
						resp,
					)
				}
			}
		})
	}
} // }}}
