// main.go
package main

import (
	"SmartWorker/SwProto"
	"encoding/binary"
	"io"
	"log"
	"net"
	"runtime"
	"strings"
	"unsafe"
)

//func init() {
//	mapIdVsHandler = make(map[uint16]IBaseHandler)
//	mapIdVsHandler[CONST_MSG_ID_MYTEST] = *(new(myproto.SMyHandler))
//}

func main() {

	go TcpService(string("0.0.0.0:9999"))

	//waiting...

	ch := make(chan int, 1)
	<-ch
}

func TcpService(strListenAddr string) {
	log.Printf("TCP: listening on %s", strListenAddr)

	tcpListener, err := net.Listen("tcp", strListenAddr)
	if err != nil {
		log.Fatalf("FATAL: listen (%s) failed - %s", strListenAddr, err.Error())
	}

	for {
		clientConn, err := tcpListener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Printf("NOTICE: temporary Accept() failure - %s", err.Error())
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("ERROR: listener.Accept() - %s", err.Error())
			}
			break
		}

		go HandleSession(clientConn)
	}

	log.Printf("TCP: closing %s", tcpListener.Addr().String())
}

//func HTTPServer(listener net.Listener) {
//	log.Printf("%s: listening on %s", proto_name, listener.Addr().String())

//	server := &http.Server{
//		Handler: handler,
//	}
//	err := server.Serve(listener)
//	// theres no direct way to detect this error because it is not exposed
//	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
//		log.Printf("ERROR: http.Serve() - %s", err.Error())
//	}

//	log.Printf("%s: closing %s", proto_name, listener.Addr().String())
//}

func HandleSession(clientConn net.Conn) {
	log.Printf("TCP: new client(%s)", clientConn.RemoteAddr())

	for {
		// The client should initialize itself by sending a 4 byte sequence indicating
		// the version of the protocol that it intends to communicate, this will allow us
		// to gracefully upgrade the protocol away from text/line oriented to whatever...
		var bytesReadBuf [SwProto.CONST_MSG_HEADER_LEN]byte
		_, err := io.ReadFull(clientConn, bytesReadBuf[:])
		if err != nil {
			log.Printf("ERROR: failed to read protocol version - %s", err.Error())
			return
		}

		msgHdr := (*SwProto.SMsgHeader)(unsafe.Pointer(&bytesReadBuf))

		if handler, OK := SwProto.GSGlobalObj.MapIdVsHandler[msgHdr.Ui16Opcode]; OK {
			handler.HandleMsg(&clientConn, msgHdr.Ui32MsgLen)

		} else {
			log.Fatalln("bad error!")
			return
		}
	}
	//int32(binary.BigEndian.Uint32(buf))

	//log.Printf("CLIENT(%s): desired protocol magic '%s'", clientConn.RemoteAddr(), protocolMagic)

	//var prot util.Protocol
	//switch protocolMagic {
	//case "  V2":
	//	prot = &protocolV2{context: p.context}
	//default:
	//	util.SendFramedResponse(clientConn, frameTypeError, []byte("E_BAD_PROTOCOL"))
	//	clientConn.Close()
	//	log.Printf("ERROR: client(%s) bad protocol magic '%s'", clientConn.RemoteAddr(), protocolMagic)
	//	return
	//}

	//err = prot.IOLoop(clientConn)
	//if err != nil {
	//	log.Printf("ERROR: client(%s) - %s", clientConn.RemoteAddr(), err.Error())
	//	return
	//}
}

// SendFramedResponse is a server side utility function to prefix data with a length header
// and frame header and write to the supplied Writer
func SendFramedResponse(w io.Writer, frameType int32, data []byte) (int, error) {
	beBuf := make([]byte, 4)
	size := uint32(len(data)) + 4

	binary.BigEndian.PutUint32(beBuf, size)
	n, err := w.Write(beBuf)
	if err != nil {
		return n, err
	}

	binary.BigEndian.PutUint32(beBuf, uint32(frameType))
	n, err = w.Write(beBuf)
	if err != nil {
		return n + 4, err
	}

	n, err = w.Write(data)
	return n + 8, err
}

//func readMPUB(r io.Reader, tmp []byte, idChan chan MessageID, maxMessageSize int64) ([]*Message, error) {
//	numMessages, err := readLen(r, tmp)
//	if err != nil {
//		return nil, util.NewFatalClientErr(err, "E_BAD_BODY", "MPUB failed to read message count")
//	}

//	if numMessages <= 0 {
//		return nil, util.NewFatalClientErr(err, "E_BAD_BODY",
//			fmt.Sprintf("MPUB invalid message count %d", numMessages))
//	}

//	messages := make([]*Message, 0, numMessages)
//	for i := int32(0); i < numMessages; i++ {
//		messageSize, err := readLen(r, tmp)
//		if err != nil {
//			return nil, util.NewFatalClientErr(err, "E_BAD_MESSAGE",
//				fmt.Sprintf("MPUB failed to read message(%d) body size", i))
//		}

//		if messageSize <= 0 {
//			return nil, util.NewFatalClientErr(nil, "E_BAD_MESSAGE",
//				fmt.Sprintf("MPUB invalid message(%d) body size %d", i, messageSize))
//		}

//		if int64(messageSize) > maxMessageSize {
//			return nil, util.NewFatalClientErr(nil, "E_BAD_MESSAGE",
//				fmt.Sprintf("MPUB message too big %d > %d", messageSize, maxMessageSize))
//		}

//		msgBody := make([]byte, messageSize)
//		_, err = io.ReadFull(r, msgBody)
//		if err != nil {
//			return nil, util.NewFatalClientErr(err, "E_BAD_MESSAGE", "MPUB failed to read message body")
//		}

//		messages = append(messages, NewMessage(<-idChan, msgBody))
//	}

//	return messages, nil
//}
