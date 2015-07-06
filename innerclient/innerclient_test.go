package innerclient

import (
	"testing"
)

/*server should like this

  func main() {
    var server innerserver.InnerServer
    server.ListenIp = ""
    server.ListenPort = 3333
    server.TransactProcess = testReq
    server.CheckConnStr = "inner conn check"

    server.Start()
    return
  }

  func testReq(reqType byte, reqData []byte) []byte {
    switch reqType {
        case 0:
        case 1:
        case 2:
        case 99:
            return []byte("--..--..")
        default:
    }
    return []byte("...")
  }

*/
func Benchmark_Request(b *testing.B) {

	var clientReq InnerClient
	clientReq.ConnPoolSize = 1
	clientReq.ServerAddr = "127.0.0.1:3333"
	clientReq.ErrRetryTimes = 3
	clientReq.CheckConnStr = "--..--.."
	clientReq.Init()

	for i := 0; i < b.N; i++ {
		res, err := clientReq.Request([]byte("hhhhhhhhhh"), 99)
		if err != nil {
			b.Error(err)
		}
		if string(res) != "----" {
			b.Error("not equal")
		}
	}
}
