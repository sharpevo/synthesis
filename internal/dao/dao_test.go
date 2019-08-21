package dao_test

import (
	"bytes"
	"encoding/binary"
	"strings"
	"synthesis/dao"
	"testing"
)

func TestByteSequnce(t *testing.T) {
	testList := []struct {
		argument  dao.Argument
		expected  []byte
		errString string
	}{
		{
			argument: dao.Argument{
				Value:      11.22, // float64 as default in golang
				ByteOrder:  binary.LittleEndian,
				ByteLength: 4,
			},
			errString: "unexpected length",
		},
		{
			argument: dao.Argument{
				Value:      float32(11.22),
				ByteOrder:  binary.LittleEndian,
				ByteLength: 4,
			},
			expected: []byte{0x1f, 0x85, 0x33, 0x41},
		},
		{
			argument: dao.Argument{
				Value:      int32(11),
				ByteOrder:  binary.LittleEndian,
				ByteLength: 4,
			},
			expected: []byte{0x0b, 0x00, 0x00, 0x00},
		},
		{
			argument: dao.Argument{
				Value:      float32(11),
				ByteOrder:  binary.LittleEndian,
				ByteLength: 4,
			},
			expected: []byte{0x00, 0x00, 0x30, 0x41},
		},
	}

	for i, test := range testList {
		actual, err := test.argument.ByteSequence()
		t.Logf("#%d", i)
		if test.errString != "" {
			if !strings.Contains(err.Error(), test.errString) {
				t.Errorf("unexpected error: %q", err.Error())
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(test.expected, actual) {
				t.Errorf(
					"\nEXPECT: '%#v'\nGET: '%#v'\n",
					test.expected,
					actual,
				)
			}
		}
	}

}

func TestNewInt32Argument(t *testing.T) {
	testList := []struct {
		errString string
		expected  interface{}
		input     interface{}
		order     binary.ByteOrder
	}{
		{
			//errString: "invalid type of argument",
			errString: "invalid ",
			input:     11.22,
			order:     binary.LittleEndian,
		},
		{
			expected: []byte{0x0b, 0x00, 0x00, 0x00},
			input:    int32(11),
			order:    binary.LittleEndian,
		},
		{
			expected: []byte{0x00, 0x0a, 0x00, 0x00},
			input:    int32(2560),
			order:    binary.LittleEndian,
		},
		{
			errString: "strconv.ParseInt",
			input:     "test",
			order:     binary.LittleEndian,
		},
	}

	for i, test := range testList {
		t.Logf("#%d", i)
		argument, err := dao.NewInt32Argument(test.input, test.order)
		if test.errString != "" {
			if !strings.Contains(err.Error(), test.errString) {
				t.Errorf("unexpected error: %q", err.Error())
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			actual, err := argument.ByteSequence()
			if err != nil || !bytes.Equal(actual, test.expected.([]byte)) {
				t.Errorf(
					"\nEXPECT: '%#v'\nGET: '%#v'\n",
					test.expected,
					actual,
				)
			}
		}
	}
}

func TestNewFloat32Argument(t *testing.T) {
	testList := []struct {
		errString string
		expected  interface{}
		input     interface{}
		order     binary.ByteOrder
	}{
		{
			errString: "invalid ",
			input:     11.22,
			order:     binary.LittleEndian,
		},
		{
			expected: []byte{0x00, 0x00, 0x30, 0x41},
			input:    float32(11),
			order:    binary.LittleEndian,
		},
		{
			expected: []byte{0x1f, 0x85, 0x33, 0x41},
			input:    float32(11.22),
			order:    binary.LittleEndian,
		},
		{
			errString: "syntax error scanning number",
			input:     "test",
			order:     binary.LittleEndian,
		},
	}

	for i, test := range testList {
		t.Logf("#%d", i)
		argument, err := dao.NewFloat32Argument(test.input, test.order)
		if test.errString != "" {
			if !strings.Contains(err.Error(), test.errString) {
				t.Errorf("unexpected error: %q", err.Error())
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			actual, err := argument.ByteSequence()
			if err != nil || !bytes.Equal(actual, test.expected.([]byte)) {
				t.Errorf(
					"\nEXPECT: '%#v'\nGET: '%#v'\n",
					test.expected,
					actual,
				)
			}
		}
	}
}

func TestNewArgument(t *testing.T) {
	testList := []struct {
		errString string
		expected  interface{}
		input     interface{}
		order     binary.ByteOrder
		length    int
	}{
		{
			errString: "invalid type of argument",
			input:     11.22,
			order:     binary.LittleEndian,
			length:    5,
		},
		{
			expected: []byte{0x00, 0x01, 0x02, 0x03, 0x04},
			input:    "0001020304",
			order:    binary.LittleEndian,
			length:   5,
		},
		{
			errString: "encoding/hex",
			expected:  []byte{0x01, 0x02},
			input:     "test",
			order:     binary.LittleEndian,
			length:    5,
		},
	}

	for i, test := range testList {
		t.Logf("#%d", i)
		argument, err := dao.NewArgument(test.input, test.order, test.length)
		if test.errString != "" && err != nil {
			//if !strings.Contains(err.Error(), test.errString) {
			//t.Errorf("unexpected error: %q", err.Error())
			//}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			actual, err := argument.ByteSequence()
			if test.errString != "" && err != nil {
				if !strings.Contains(err.Error(), test.errString) {
					t.Errorf("unexpected error: %q", err.Error())
				}
			} else if !bytes.Equal(actual, test.expected.([]byte)) {
				t.Errorf(
					"\nEXPECT: '%#v'\nGET: '%#v'\n",
					test.expected,
					actual,
				)
			}
		}
	}
}
