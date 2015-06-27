# innerconn

#How to use:

go get github.com/conc/innerconn/innerserver   
go get github.com/conc/innerconn/innerclient   


#Examples:

```golang
client.go:
package main

import (
	"github.com/conc/innerconn/innerclient"
	"log"
)

func main() {
	var clientReq innerclient.InnerClient
	clientReq.ConnPoolSize = 10
	clientReq.ServerAddr = "127.0.0.1:3333"
	clientReq.ErrRetryTimes = 3
	clientReq.CheckConnStr = "--..--.."
	clientReq.Init()

	for i := 0; i < 1000000; i++ {
		ret, err := clientReq.Request([]byte("hhhhhhhhhh"), 99)
		log.Println(i, string(ret))
		if err != nil {
			log.Println(err)
		}
	}

	return
}

server.go:
package main

import (
	"github.com/conc/innerconn/innerserver"
	"log"
)

func main() {
	var server innerserver.InnerServer
	server.ListenIp = ""
	server.ListenPort = 3333
	server.TransactProcess = dealTestReq
	server.CheckConnStr = "--..--.."

	server.Start()
	return
}

func dealTestReq(reqType byte, reqData []byte) []byte {

	switch reqType {
	case 0:
	case 1:
	case 2:
	case 99:
		log.Println(string(reqData))
		return []byte("----")
	default:
	}

	return []byte("...")
}
```


