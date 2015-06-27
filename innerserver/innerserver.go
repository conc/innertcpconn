package innerserver

import (
	"log"
	"net"
	"os"
	"strings"
)

type TransactBussiness func(byte, []byte) []byte

type InnerServer struct {
	ListenIp        string
	ListenPort      int
	TransactProcess TransactBussiness
	CheckConnStr    string
}

func (c *InnerServer) Start() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(c.ListenIp), c.ListenPort, ""})
	if err != nil {
		log.Println("err ListenTCP")
		os.Exit(1)
		return
	}

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("Error AcceptTCP:", err)
			continue
		}
		if !c.checkConn(conn) {
			conn.Close()
			continue
		}
		go c.dealConn(conn)
	}
}

func (c *InnerServer) checkConn(conn *net.TCPConn) bool {
	dataBuf := make([]byte, len(c.CheckConnStr))
	readLen, err := (*conn).Read(dataBuf[:])
	if err != nil {
		return false
	}
	if !strings.EqualFold(string(dataBuf[0:readLen]), c.CheckConnStr) {
		return false
	}

	(*conn).Write([]byte(c.CheckConnStr))
	return true
}

func (c *InnerServer) dealConn(conn *net.TCPConn) {
	defer conn.Close()
	data := make([]byte, 1024)
	for {
		datalen, err := conn.Read(data)
		if err != nil {
			log.Println("Error conn.Read:", err)
			break
		}

		res := c.dealReceiveData(data[0:datalen])
		resByte := connStuToBytes(res)
		conn.Write(resByte)
	}
	return
}

func (c *InnerServer) dealReceiveData(data []byte) *connStu {
	var res connStu

	req, err := bytesToConnStu(data)
	if err != nil {
		return &res
	}

	res.Data = c.TransactProcess(req.RequestType, req.Data)
	res.RequestType = req.RequestType
	res.RequestId = req.RequestId
	res.DataLen = uint64(len(res.Data))
	return &res
}
