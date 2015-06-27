package innerserver

import (
	"reflect"
	"testing"
)

func Test_BytesUintExchange(t *testing.T) {
	var cc uint64 = 7654321
	ccBytes := uint64ToBytes(cc)
	res := bytesToUint64(ccBytes)
	if res != cc {
		t.Error("err change")
	}
}

func Test_StuBytesExchange(t *testing.T) {
	var cc connStu
	cc.RequestId = 12345
	cc.RequestType = 99
	cc.Data = []byte("hehehehehe")
	cc.DataLen = uint64(len(cc.Data))

	byteData := connStuToBytes(&cc)
	res, err := bytesToConnStu(byteData)
	if err != nil {
		t.Error("err change", err)
	}

	if !reflect.DeepEqual(cc, *res) {
		t.Error("err change", err)
	}
}
