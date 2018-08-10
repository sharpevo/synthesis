package ricoh_g5

import (
	"encoding/binary"
	"log"
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
	bitsPerPixel string,
	width string,
	lineBufferSize string,
	lineBuffer string,
) (resp string, err error) {

	bitsPerPixelBytes, err := Int32ByteSequence(bitsPerPixel)
	if err != nil {
		return resp, err
	}
	widthBytes, err := Int32ByteSequence(width)
	if err != nil {
		return resp, err
	}
	lineBufferSizeArgument, err := dao.NewInt32Argument(lineBufferSize, binary.LittleEndian)
	if err != nil {
		return resp, err
	}
	lineBufferSizeBytes, err := lineBufferSizeArgument.ByteSequence()
	if err != nil {
		return resp, err
	}
	length := lineBufferSizeArgument.Value.(int32)
	lineBufferBytes, err := VariantByteSequence(lineBuffer, int(length))
	if err != nil {
		return resp, err
	}

	req := PrintDataUnit.Request()
	reqBytes := req.Bytes()
	reqBytes = append(reqBytes, bitsPerPixelBytes...)
	reqBytes = append(reqBytes, widthBytes...)
	reqBytes = append(reqBytes, lineBufferSizeBytes...)
	reqBytes = append(reqBytes, lineBufferBytes...)
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

func (d *Dao) SendWaveform(
	headBoardIndex string,
	rowIndexOfHeadBoard string,
	voltagePercentage string,
	segmentCount string,
	segment string,
) (resp string, err error) {
	headBoardIndexBytes, err := Int32ByteSequence(headBoardIndex)
	if err != nil {
		return resp, err
	}
	rowIndexOfHeadBoardBytes, err := Int32ByteSequence(rowIndexOfHeadBoard)
	if err != nil {
		return resp, err
	}
	voltagePercentageBytes, err := Float32ByteSequence(voltagePercentage)
	if err != nil {
		return resp, err
	}
	segmentCountArgument, err := dao.NewInt32Argument(segmentCount, binary.LittleEndian)
	if err != nil {
		return resp, err
	}
	segmentCountBytes, err := segmentCountArgument.ByteSequence()
	if err != nil {
		return resp, err
	}
	length := segmentCountArgument.Value.(int32)
	segmentBytes, err := VariantByteSequence(segment, int(length))
	if err != nil {
		return resp, err
	}

	req := WaveformUnit.Request()
	reqBytes := req.Bytes()
	reqBytes = append(reqBytes, headBoardIndexBytes...)
	reqBytes = append(reqBytes, rowIndexOfHeadBoardBytes...)
	reqBytes = append(reqBytes, voltagePercentageBytes...)
	reqBytes = append(reqBytes, segmentCountBytes...)
	reqBytes = append(reqBytes, segmentBytes...)

	respBytes, err := d.TCPClient.Send(
		reqBytes,
		WaveformUnit.ComResp(),
	)
	resp = string(respBytes)
	if err != nil {
		log.Println("ERR:", err)
		return resp, err
	}
	return resp, nil
}

func Int32ByteSequence(input interface{}) (output []byte, err error) {
	argument, err := dao.NewInt32Argument(input, binary.LittleEndian)
	if err != nil {
		return output, err
	}
	output, err = argument.ByteSequence()
	if err != nil {
		return output, err
	}
	return
}

func Float32ByteSequence(input interface{}) (output []byte, err error) {
	argument, err := dao.NewFloat32Argument(input, binary.LittleEndian)
	if err != nil {
		return output, err
	}
	output, err = argument.ByteSequence()
	if err != nil {
		return output, err
	}
	return
}

func VariantByteSequence(input interface{}, length int) (output []byte, err error) {
	argument, err := dao.NewArgument(input, binary.LittleEndian, length)
	if err != nil {
		return output, err
	}
	output, err = argument.ByteSequence()
	if err != nil {
		return output, err
	}
	return
}
