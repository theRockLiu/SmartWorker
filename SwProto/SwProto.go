// SwProto.go
package SwProto

import (
	proto "code.google.com/p/goprotobuf/proto"
	"fmt"
	"io"
	"log"
	"net"
)

type IBaseHandler interface {
	HandleMsg(*net.Conn, uint32) error
}

type SGlobalObj struct {
	mapIdVsHandler map[uint16]IBaseHandler
}

var gSGlobalObj SGlobalObj

type SMsgHeader struct {
	ui32MsgLen uint32
	ui16Opcode uint16
}

var tmptmptmptmp SMsgHeader

const CONST_MSG_HEADER_LEN = unsafe.Sizeof(tmptmptmptmp)

const (
	CST_MSGID_REG      = 1
	CST_MSGID_RES      = 2
	CST_MSGID_SET_NAME = 3
	CST_MSGID_GET_NAME = 4
)

////////////////////////////////////////////////////////////////////////////////

type SRegHandler struct {
}

func (this SRegHandler) HandleMsg(conn *net.Conn, ui32BodyLen uint32) error {
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

	fmt.Println(regreq.GetLabel(), regreq.GetType())

	return nil
}

type SSetNameHandler struct {
}

func (this SSetNameHandler) HandleMsg(conn *net.Conn, ui32BodyLen uint32) error {

}

type SGetNameHandler struct {
}

func (this SGetNameHandler) HandleMsg(conn *net.Conn, ui32BodyLen uint32) error {

}

////////////////////////////////////////////////////////////////////////////////
func init() {
	gSGlobalObj.mapIdVsHandler[CST_MSGID_REG] = SRegHandler{}

}
