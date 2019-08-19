package ricoh_g5

import (
	"posam/dao"
)

var ErrorCodeUnit,
	PrinterStatusUnit,
	PrintDataUnit,
	WaveformUnit dao.Unit

func init() {

	// ErrorCode{{{
	ErrorCodeRequest := dao.Request{}
	ErrorCodeRequest.Function = 0x01
	ErrorCodeRequest.Arguments = []byte{
		0x00,
		0x00,
		0x00,
	}
	ErrorCodeRecResp := []byte{}
	ErrorCodeComResp := []byte{
		0x01,
		0x00,
		0x00,
		0x00,

		0x00,
		0x00,
		0x00,
		0x00,
	}
	ErrorCodeUnit.SetRequest(ErrorCodeRequest)
	ErrorCodeUnit.SetRecResp(ErrorCodeRecResp)
	ErrorCodeUnit.SetComResp(ErrorCodeComResp)
	// }}}

	// PrinterStatus{{{
	PrinterStatusRequest := dao.Request{}
	PrinterStatusRequest.Function = 0x02
	PrinterStatusRequest.Arguments = []byte{
		0x00,
		0x00,
		0x00,
	}
	PrinterStatusRecResp := []byte{}
	PrinterStatusComResp := []byte{
		0x02,
		0x00,
		0x00,
		0x00,

		0x03,
		0x00,
		0x00,
		0x00,
	}
	PrinterStatusUnit.SetRequest(PrinterStatusRequest)
	PrinterStatusUnit.SetRecResp(PrinterStatusRecResp)
	PrinterStatusUnit.SetComResp(PrinterStatusComResp)
	// }}}

	// PrintData{{{
	PrintDataRequest := dao.Request{}
	PrintDataRequest.Function = 0x03
	PrintDataRequest.Arguments = []byte{
		0x00,
		0x00,
		0x00,
	}
	PrintDataRecResp := []byte{}
	PrintDataComResp := []byte{
		0x03,
		0x00,
		0x00,
		0x00,

		0x00,
		0x00,
		0x00,
		0x00,
	}
	PrintDataUnit.SetRequest(PrintDataRequest)
	PrintDataUnit.SetRecResp(PrintDataRecResp)
	PrintDataUnit.SetComResp(PrintDataComResp)
	// }}}

	// WaveForm{{{
	WaveformRequest := dao.Request{}
	WaveformRequest.Function = 0x04
	WaveformRequest.Arguments = []byte{
		0x00,
		0x00,
		0x00,
	}
	WaveformRecResp := []byte{}
	WaveformComResp := []byte{
		0x04,
		0x00,
		0x00,
		0x00,

		0x00,
		0x00,
		0x00,
		0x00,
	}
	WaveformUnit.SetRequest(WaveformRequest)
	WaveformUnit.SetRecResp(WaveformRecResp)
	WaveformUnit.SetComResp(WaveformComResp)
	// }}}

}
