package ricoh_g5

import (
	"log"
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
