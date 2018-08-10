package dao

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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
	var buf = new(bytes.Buffer)
	err = binary.Write(buf, a.ByteOrder, a.Value)
	if err != nil {
		return output, err
	}
	output = buf.Bytes()
	if len(output) != a.ByteLength {
		return output, fmt.Errorf(
			"%v is translated with unexpected length",
			a.Value)
	}
	return output, nil
}
