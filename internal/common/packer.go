package common

import (
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/util"
)

func PackUserReq(svc kiwi.TSvc, code kiwi.TCode, seqId uint32, payload []byte) ([]byte, *util.Err) {
	var buffer util.ByteBuffer
	buffer.InitCap(64)
	sc := kiwi.MergeSvcCode(svc, code)
	buffer.WUint16(sc)
	buffer.WUint32(seqId)
	_, e := buffer.Write(payload)
	if e != nil {
		return nil, util.NewErr(util.EcWriteFail, util.M{
			"error": e,
		})
	}
	return buffer.All(), nil
}

func UnpackUserReq(bytes []byte) (svc kiwi.TSvc, code kiwi.TCode, seqId uint32, payload []byte, err *util.Err) {
	var buffer util.ByteBuffer
	buffer.InitBytes(bytes)
	sc, err := buffer.RUint16()
	if err != nil {
		return
	}
	seqId, err = buffer.RUint32()
	if err != nil {
		return
	}
	svc, code = kiwi.SplitSvcCode(sc)
	payload = buffer.RAvailable()
	return
}

func PackUserOk(seqId uint32, payload []byte) ([]byte, *util.Err) {
	var buffer util.ByteBuffer
	buffer.InitCap(5 + uint32(len(payload)))
	buffer.WUint8(0)
	buffer.WUint32(seqId)
	_, e := buffer.Write(payload)
	if e != nil {
		return nil, util.NewErr(util.EcWriteFail, util.M{
			"error": e,
		})
	}
	return buffer.All(), nil
}

func UnpackUserResOk(bytes []byte) (seqId uint32, payload []byte, err *util.Err) {
	var buffer util.ByteBuffer
	buffer.InitBytes(bytes)
	buffer.SetPos(1)
	seqId, err = buffer.RUint32()
	if err != nil {
		return
	}
	payload = buffer.RAvailable()
	return
}

func PackUserFail(seqId uint32, resCode uint16) ([]byte, *util.Err) {
	var buffer util.ByteBuffer
	buffer.InitCap(7)
	buffer.WUint8(1)
	buffer.WUint32(seqId)
	buffer.WUint16(resCode)
	return buffer.All(), nil
}

func UnpackUserResFail(bytes []byte) (seqId uint32, resCode uint16, err *util.Err) {
	var buffer util.ByteBuffer
	buffer.InitBytes(bytes)
	buffer.SetPos(1)
	seqId, err = buffer.RUint32()
	if err != nil {
		return
	}
	resCode, err = buffer.RUint16()
	return
}

func PackUserPus(svc kiwi.TSvc, code kiwi.TCode, ntc []byte) ([]byte, *util.Err) {
	var buffer util.ByteBuffer
	buffer.InitCap(3 + uint32(len(ntc)))
	buffer.WUint8(2)
	buffer.WUint16(kiwi.MergeSvcCode(svc, code))
	_, e := buffer.Write(ntc)
	if e != nil {
		return nil, util.NewErr(util.EcWriteFail, util.M{
			"error": e,
		})
	}
	return buffer.All(), nil
}

func UnpackUserPus(bytes []byte) (svc kiwi.TSvc, code kiwi.TCode, payload []byte, err *util.Err) {
	var buffer util.ByteBuffer
	buffer.InitBytes(bytes)
	buffer.SetPos(1)
	sc, err := buffer.RUint16()
	if err != nil {
		return
	}
	svc, code = kiwi.SplitSvcCode(sc)
	payload = buffer.RAvailable()
	return
}
