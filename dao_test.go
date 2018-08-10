package dao_test

import (
	"bytes"
	"encoding/binary"
	"posam/dao"
	"strings"
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
			expected:  []byte("test"),
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
