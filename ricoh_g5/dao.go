package ricoh_g5

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"posam/dao"
	"posam/protocol/tcp"
	"posam/util/concurrentmap"
)

var deviceMap *concurrentmap.ConcurrentMap

func init() {
	deviceMap = concurrentmap.NewConcurrentMap()
}

type Dao struct {
	DeviceAddress string
	TCPClient     tcp.TCPClienter
}

func AddInstance(dao *Dao) {
	deviceMap.Set(dao.DeviceAddress, dao)
}

func Instance(address string) *Dao {
	if device, ok := deviceMap.Get(address); ok {
		return device.(*Dao)
	} else {
		return nil
	}
}

func (d *Dao) QueryErrorCode() (resp string, err error) {
	req := ErrorCodeUnit.Request()
	respBytes, err := d.TCPClient.Send(
		req.Bytes(),
		ErrorCodeUnit.ComResp(),
	)
	resp = string(respBytes)
	if err != nil {
		log.Println("ERR:", err)
		return resp, err
	}
	return resp, nil
}

func (d *Dao) QueryPrinterStatus() (resp string, err error) {
	req := PrinterStatusUnit.Request()
	respBytes, err := d.TCPClient.Send(
		req.Bytes(),
		PrinterStatusUnit.ComResp(),
	)
	resp = string(respBytes)
	if err != nil {
		log.Println("ERR:", err)
		return resp, err
	}
	return resp, nil
}

func (d *Dao) PrintData(
	bitsPerPixel []byte,
	width []byte,
	lineBufferSize []byte,
	lineBuffer []byte,
) (resp string, err error) {
	req := PrintDataUnit.Request()
	reqBytes := req.Bytes()
	reqBytes = append(reqBytes, bitsPerPixel...)
	reqBytes = append(reqBytes, width...)
	reqBytes = append(reqBytes, lineBufferSize...)
	reqBytes = append(reqBytes, lineBuffer...)
	respBytes, err := d.TCPClient.Send(
		reqBytes,
		PrintDataUnit.ComResp(),
	)
	resp = string(respBytes)
	if err != nil {
		log.Println("ERR:", err)
		return resp, err
	}
	return resp, nil
}

func NewArgument(input interface{}) (argument dao.Argument, err error) {
	argument.ByteOrder = binary.LittleEndian
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
		argument.Value = f
		return argument, err
	case int32, uint32, float32:
		argument.Value = v
		return
	default:
		return argument, fmt.Errorf("invalid type of argument: %s", v)
	}
}
