// SwProto.go
package SwProto

import (
	"SmartWorker/msgs"
	proto "code.google.com/p/goprotobuf/proto"
	//	"fmt"
	"io"
	"log"
	"net"
	"unsafe"
)

type IBaseHandler interface {
	HandleMsg(*net.Conn, uint32) error
}

type SGlobalObj struct {
	MapIdVsHandler map[uint16]IBaseHandler
}

var GSGlobalObj SGlobalObj

type SMsgHeader struct {
	Ui32MsgLen uint32
	Ui16Opcode uint16
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

	regreq := &msgs.Regreq{}
	//myPerson := &proto.Person{}
	err = proto.Unmarshal(bytesReadBuf, regreq)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return err
	}

	//fmt.Println(regreq.GetLabel(), regreq.GetType())

	return nil
}

type SSetNameHandler struct {
}

func (this SSetNameHandler) HandleMsg(conn *net.Conn, ui32BodyLen uint32) error {
	log.Println("handling set name msg")

	return nil
}

type SGetNameHandler struct {
}

func (this SGetNameHandler) HandleMsg(conn *net.Conn, ui32BodyLen uint32) error {
	log.Println("handliing get name msg")

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func init() {
	GSGlobalObj.MapIdVsHandler = make(map[uint16]IBaseHandler)
	GSGlobalObj.MapIdVsHandler[CST_MSGID_REG] = SRegHandler{}
	GSGlobalObj.MapIdVsHandler[CST_MSGID_SET_NAME] = SSetNameHandler{}
	GSGlobalObj.MapIdVsHandler[CST_MSGID_GET_NAME] = SGetNameHandler{}

}
