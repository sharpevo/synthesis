package dao

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"posam/interpreter"
	"strconv"
)

const NAME = "UNKNOWN"

var InstructionMap interpreter.InstructionMapt

func init() {
	InstructionMap = make(interpreter.InstructionMapt)
}

type Frame struct {
	Address   byte
	Function  byte
	Arguments []byte
}

type Bit struct {
	Frame
	CRC []byte
}

type Request struct {
	Bit
}

type ReceivedResponse []byte
type CompletedResponse []byte

type Unit struct {
	request Request
	recResp ReceivedResponse
	comResp CompletedResponse
}

func (u *Unit) Request() Request {
	return u.request
}

func (u *Unit) RecResp() ReceivedResponse {
	return u.recResp
}

func (u *Unit) ComResp() CompletedResponse {
	return u.comResp
}

func (u *Unit) SetRequest(request Request) {
	u.request = request
}

func (u *Unit) SetRecResp(recResp ReceivedResponse) {
	u.recResp = recResp
}

func (u *Unit) SetComResp(comResp CompletedResponse) {
	u.comResp = comResp
}

func (r *Request) Bytes() (output []byte) {
	if r.Address != 0 {
		output = append(output, r.Address)
	}
	output = append(output, r.Function)
	output = append(output, r.Arguments...)
	output = append(output, r.CRC...)
	return
}

type Argument struct {
	Value      interface{}
	ByteOrder  binary.ByteOrder
	ByteLength int
}

func (a *Argument) ByteSequence() (output []byte, err error) {
	switch v := a.Value.(type) {
	case string:
		output, err = hex.DecodeString(v)
		if err != nil {
			return output, err
		}
	default:
		var buf = new(bytes.Buffer)
		err = binary.Write(buf, a.ByteOrder, a.Value)
		if err != nil {
			return output, err
		}
		output = buf.Bytes()
	}
	if len(output) != a.ByteLength {
		fmt.Println("Error:", a.Value, output, a.ByteLength)
		return output, fmt.Errorf(
			"%v is translated with unexpected length %d (%d)",
			a.Value,
			len(output),
			a.ByteLength,
		)
	}
	return output, nil
}

func NewArgument(
	input interface{},
	order binary.ByteOrder,
	length int,
) (argument Argument, err error) {
	argument.ByteOrder = order
	argument.ByteLength = length
	switch v := input.(type) {
	case string:
		argument.Value = v
		return
	default:
		return argument, fmt.Errorf("invalid variant argument: %v", v)
	}

}

func NewInt32Argument(
	input interface{},
	order binary.ByteOrder,
) (argument Argument, err error) {
	argument.ByteOrder = order
	argument.ByteLength = 4
	switch v := input.(type) {
	case string:
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return argument, err
		}
		return NewInt32Argument(i, order)
	case int32:
		argument.Value = v
		return
	case int64:
		argument.Value = int32(v)
		return
	default:
		return argument, fmt.Errorf("invalid int32 argument: %v", v)
	}
}

func NewFloat32Argument(
	input interface{},
	order binary.ByteOrder,
) (argument Argument, err error) {
	argument.ByteOrder = order
	argument.ByteLength = 4
	switch v := input.(type) {
	case string:
		bigF, _, err := big.ParseFloat(v, 10, 24, big.ToNearestEven)
		if err != nil {
			return argument, err
		}
		f, acc := bigF.Float32()
		if acc != big.Exact {
			return argument, fmt.Errorf("failed to parse '%v' exactly", input)
		}
		return NewFloat32Argument(f, order)
	case float32:
		argument.Value = v
		return
	default:
		return argument, fmt.Errorf("invalid float32 argument: %v", v)
	}
}
