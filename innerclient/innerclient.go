package innerclient

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type InnerClient struct {
	ConnPoolSize  int
	connPool      chan *net.Conn
	ServerAddr    string
	ErrRetryTimes int
	CheckConnStr  string
}

var reqNum uint64

func (c *InnerClient) Init() {
	log.Println("Loading conn...............")
	c.connPool = make(chan *net.Conn, c.ConnPoolSize)
	for i := 0; i < c.ConnPoolSize; i++ {
		connTemp, err := c.createConn()
		if err != nil {
			log.Println("get conn error!")
			os.Exit(1)
		}
		c.connPool <- connTemp
	}

	log.Println("Loading conn over..........")
	return
}

func (c *InnerClient) Request(reqData []byte, reqType byte) (res []byte, err error) {

	req := connStu{RequestType: reqType, Data: reqData}
	req.DataLen = uint64(len(reqData))
	req.RequestId = atomic.AddUint64(&reqNum, 1)

	var ret *connStu
	for i := 0; i < c.ErrRetryTimes+1; i++ {
		ret, err = c.sendReceive(&req)
		if err == nil {
			res = ret.Data
			return
		} else {
			continue
		}
	}
	return
}

func (c *InnerClient) sendReceive(req *connStu) (res *connStu, err error) {
	var clientConn *net.Conn
	select {
	case clientConn = <-c.connPool:
	case <-time.After(2 * time.Second):
		err = errors.New("Get nothing from connPool")
		return
	}

	err = c.send(clientConn, req)
	if err != nil {
		return
	}

	res, err = c.receive(clientConn)
	if err != nil {
		return
	}

	if res.RequestId != req.RequestId {
		err = errors.New("receive data error(requsetid error)!")
		return
	}
	return
}

func (c *InnerClient) send(clientConn *net.Conn, req *connStu) (err error) {
	reqByte := connStuToBytes(req)
	(*clientConn).SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = (*clientConn).Write(reqByte)
	if err != nil {
		go c.dealErrorConn(clientConn)
		return
	}
	return
}

func (c *InnerClient) receive(clientConn *net.Conn) (ret *connStu, err error) {
	var dataBuf = make([]byte, 1024)
	var readLen int
	(*clientConn).SetReadDeadline(time.Now().Add(10 * time.Second))
	readLen, err = (*clientConn).Read(dataBuf[:])
	if err != nil {
		go c.dealErrorConn(clientConn)
		return
	}
	ret, err = bytesToConnStu(dataBuf[0:readLen])
	c.connPool <- clientConn
	return
}

func (c *InnerClient) createConn() (*net.Conn, error) {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		return nil, err
	}

	if !c.checkConn(&conn) {
		conn.Close()
		return nil, errors.New("checkConn error")
	}

	return &conn, nil
}

func (c *InnerClient) checkConn(conn *net.Conn) bool {
	(*conn).Write([]byte(c.CheckConnStr))
	dataBuf := make([]byte, len(c.CheckConnStr))
	readLen, err := (*conn).Read(dataBuf[:])
	if err != nil {
		return false
	}
	if !strings.EqualFold(string(dataBuf[0:readLen]), c.CheckConnStr) {
		return false
	}
	return true
}

func (c *InnerClient) dealErrorConn(clientConn *net.Conn) {
	(*clientConn).Close()
	var err error
	for {
		if clientConn, err = c.createConn(); err != nil {
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	c.connPool <- clientConn
	return
}
