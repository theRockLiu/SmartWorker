// mytest.go
package myproto

import (
	proto "code.google.com/p/goprotobuf/proto"
	"io"
	"log"
	"net"
)

type SMyHandler struct {
	i int
}

func (this SMyHandler) HandleMsg(conn *net.Conn, ui32BodyLen uint32) error {
	log.Println("SMyHandler HandleMsg: ", ui32BodyLen)

	bytesReadBuf := make([]byte, ui32BodyLen)
	_, err := io.ReadFull(*conn, bytesReadBuf)
	if err != nil {
		log.Printf("ERROR: failed to read protocol version - %s", err.Error())
		return err
	}

	regreq := &RegReq{}
	//myPerson := &proto.Person{}
	err = proto.Unmarshal(bytesReadBuf, regreq)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return err
	}

	return nil
}
