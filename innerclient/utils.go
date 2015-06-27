package innerclient

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type connStu struct {
	RequestId   uint64
	DataLen     uint64
	RequestType byte
	Data        []byte
}

func connStuToBytes(req *connStu) []byte {
	var result []byte
	result = append(result, uint64ToBytes(req.RequestId)...)
	result = append(result, uint64ToBytes(req.DataLen)...)
	result = append(result, req.RequestType)
	result = append(result, req.Data...)
	return result
}

func bytesToConnStu(data []byte) (*connStu, error) {
	var result connStu
	if len(data) < 17 {
		return &result, errors.New("receive data error,data < 17")
	}
	result.RequestId = bytesToUint64(data[0:8])
	result.DataLen = bytesToUint64(data[8:16])
	result.RequestType = data[16]
	if len(data) < int(17+result.DataLen) {
		result.RequestId = 0
		return &result, errors.New("receive data error,data < 17 +")
	}

	result.Data = data[17 : 17+result.DataLen]
	return &result, nil
}

func uint64ToBytes(data uint64) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, data)
	return b_buf.Bytes()
}

func bytesToUint64(data []byte) uint64 {
	var result uint64
	b_buf := bytes.NewBuffer(data)
	binary.Read(b_buf, binary.BigEndian, &result)
	return result
}
