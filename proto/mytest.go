// mytest.go
package proto

import (
	"SmartWorker/proto"
	"log"
	"net"
)

type SMyHandler struct {
	i int
}

func (this SMyHandler) HandleMsg(conn *net.Conn, ui32BodyLen uint32) error {
	log.Println("SMyHandler HandleMsg: ", ui32BodyLen)

	bytesReadBuf := make([]byte, ui32BodyLen)
	_, err := io.ReadFull(conn, bytesReadBuf)
	if err != nil {
		log.Printf("ERROR: failed to read protocol version - %s", err.Error())
		return err
	}

	myPerson := &proto.Person{}
	err = proto.Unmarshal(bytesReadBuf, myPerson)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return err
	}

	return nil
}
