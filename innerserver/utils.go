package innerserver

import (
	"errors"
	"unsafe"
)

var sizeUint64 = int(unsafe.Sizeof(uint64(0)))

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

// encodes a little-endian uint64 into a byte slice.
func uint64ToBytes(n uint64) []byte {
	buf := make([]byte, 8)
	for i := uint(0); i < uint(sizeUint64); i++ {
		buf[i] = byte(n >> (i * 8))
	}
	return buf
}

// decodes a little-endian uint64 from a byte slice.
func bytesToUint64(buf []byte) (n uint64) {
	n |= uint64(buf[0])
	n |= uint64(buf[1]) << 8
	n |= uint64(buf[2]) << 16
	n |= uint64(buf[3]) << 24
	n |= uint64(buf[4]) << 32
	n |= uint64(buf[5]) << 40
	n |= uint64(buf[6]) << 48
	n |= uint64(buf[7]) << 56
	return
}
