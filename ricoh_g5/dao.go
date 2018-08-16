package ricoh_g5

import (
	"encoding/binary"
	"fmt"
	"log"
	"posam/dao"
	"posam/protocol/tcp"
	"posam/util/concurrentmap"
)

var deviceMap *concurrentmap.ConcurrentMap

func init() {
	ResetInstance()
}

type Dao struct {
	DeviceAddress string
	TCPClient     tcp.Clienter
}

func AddInstance(dao *Dao) {
	deviceMap.Set(dao.DeviceAddress, dao)
}

func ResetInstance() {
	deviceMap = concurrentmap.NewConcurrentMap()
	tcp.ResetInstance()
}

func Instance(address string) *Dao {
	if address == "" {
		for device := range deviceMap.Iter() {
			return device.Value.(*Dao)
		}
	} else {
		if device, ok := deviceMap.Get(address); ok {
			return device.(*Dao)
		} else {
			return nil
		}
	}
	log.Println("nil device instance")
	return nil
}

func (d *Dao) QueryErrorCode() (resp interface{}, err error) {
	req := ErrorCodeUnit.Request()
	resp, err = d.TCPClient.Send(
		req.Bytes(),
		ErrorCodeUnit.ComResp(),
	)
	if err != nil {
		log.Println(err)
		return resp, err
	}
	return resp, nil
}

func (d *Dao) QueryPrinterStatus() (resp interface{}, err error) {
	req := PrinterStatusUnit.Request()
	resp, err = d.TCPClient.Send(
		req.Bytes(),
		PrinterStatusUnit.ComResp(),
	)
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
) (resp interface{}, err error) {

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
	resp, err = d.TCPClient.Send(
		reqBytes,
		PrintDataUnit.ComResp(),
	)
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
	segmentArgumentList []string,
) (resp interface{}, err error) {
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

	segmentList, err := Segmentify(segmentArgumentList, 13)
	if err != nil {
		return resp, err
	}

	segmentBytes, err := SegmentBytes(segmentList, length)

	actual := len(segmentBytes)
	expect := int(52 * length)
	if actual != expect {
		return resp, fmt.Errorf("%v is translated with unexpected length %d (%d)",
			segmentArgumentList, actual, expect,
		)
	}

	req := WaveformUnit.Request()
	reqBytes := req.Bytes()
	reqBytes = append(reqBytes, headBoardIndexBytes...)
	reqBytes = append(reqBytes, rowIndexOfHeadBoardBytes...)
	reqBytes = append(reqBytes, voltagePercentageBytes...)
	reqBytes = append(reqBytes, segmentCountBytes...)
	reqBytes = append(reqBytes, segmentBytes...)

	resp, err = d.TCPClient.Send(
		reqBytes,
		WaveformUnit.ComResp(),
	)
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

func Segmentify(segmentArgumentList []string, length int) (segmentList [][]string, err error) {
	if len(segmentArgumentList)%length != 0 {
		return segmentList, fmt.Errorf("invalid segment")
	}
	segment := []string{}
	for len(segmentArgumentList) >= length {
		segment, segmentArgumentList = segmentArgumentList[:length], segmentArgumentList[length:]
		segmentList = append(segmentList, segment)
	}
	return segmentList, nil
}

func SegmentBytes(segmentList [][]string, length int32) (segmentsBytes []byte, err error) {
	for _, itemList := range segmentList {
		for k, item := range itemList {
			itemBytes := []byte{}
			switch k {
			case 12:
				itemBytes, err = Int32ByteSequence(item)
				if err != nil {
					return segmentsBytes, err
				}
			default:
				itemBytes, err = Float32ByteSequence(item)
				if err != nil {
					return segmentsBytes, err
				}
			}
			segmentsBytes = append(segmentsBytes, itemBytes...)
		}
	}
	return segmentsBytes, nil
}
