// main.go
package main

import (
	"fmt"
	"net"
)

type IBaseHandler interface {
	Handle(*net.Conn)
}

var mapIdVsHandler map[int32]IBaseHandler

func main() {
	var httpListener net.Listener
	//var httpsListener net.Listener
	mapIdVsHandler = make(map[int32]IBaseHandler)

	context := &context{n}

	tcpListener, err := net.Listen("tcp", n.tcpAddr.String())
	if err != nil {
		log.Fatalf("FATAL: listen (%s) failed - %s", n.tcpAddr, err.Error())
	}
	n.tcpListener = tcpListener
	tcpServer := &tcpServer{context: context}
	n.waitGroup.Wrap(func() { util.TCPServer(n.tcpListener, tcpServer) })

	//if n.tlsConfig != nil && n.httpsAddr != nil {
	//	httpsListener, err = tls.Listen("tcp", n.httpsAddr.String(), n.tlsConfig)
	//	if err != nil {
	//		log.Fatalf("FATAL: listen (%s) failed - %s", n.httpsAddr, err.Error())
	//	}
	//	n.httpsListener = httpsListener
	//	httpsServer := &httpServer{
	//		context:     context,
	//		tlsEnabled:  true,
	//		tlsRequired: true,
	//	}
	//	n.waitGroup.Wrap(func() { util.HTTPServer(n.httpsListener, httpsServer, "HTTPS") })
	//}

	httpListener, err = net.Listen("tcp", n.httpAddr.String())
	if err != nil {
		log.Fatalf("FATAL: listen (%s) failed - %s", n.httpAddr, err.Error())
	}
	n.httpListener = httpListener
	httpServer := &httpServer{
		context:     context,
		tlsEnabled:  false,
		tlsRequired: n.options.TLSRequired,
	}
	n.waitGroup.Wrap(func() { util.HTTPServer(n.httpListener, httpServer, "HTTP") })
}

func TCPServer(listener net.Listener) {
	log.Printf("TCP: listening on %s", listener.Addr().String())

	for {
		clientConn, err := listener.Accept()
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

	log.Printf("TCP: closing %s", listener.Addr().String())
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

	// The client should initialize itself by sending a 4 byte sequence indicating
	// the version of the protocol that it intends to communicate, this will allow us
	// to gracefully upgrade the protocol away from text/line oriented to whatever...
	buf := make([]byte, 4)
	_, err := io.ReadFull(clientConn, buf)
	if err != nil {
		log.Printf("ERROR: failed to read protocol version - %s", err.Error())
		return
	}
	int32(binary.BigEndian.Uint32(buf))

	log.Printf("CLIENT(%s): desired protocol magic '%s'", clientConn.RemoteAddr(), protocolMagic)

	var prot util.Protocol
	switch protocolMagic {
	case "  V2":
		prot = &protocolV2{context: p.context}
	default:
		util.SendFramedResponse(clientConn, frameTypeError, []byte("E_BAD_PROTOCOL"))
		clientConn.Close()
		log.Printf("ERROR: client(%s) bad protocol magic '%s'", clientConn.RemoteAddr(), protocolMagic)
		return
	}

	err = prot.IOLoop(clientConn)
	if err != nil {
		log.Printf("ERROR: client(%s) - %s", clientConn.RemoteAddr(), err.Error())
		return
	}
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

func readMPUB(r io.Reader, tmp []byte, idChan chan MessageID, maxMessageSize int64) ([]*Message, error) {
	numMessages, err := readLen(r, tmp)
	if err != nil {
		return nil, util.NewFatalClientErr(err, "E_BAD_BODY", "MPUB failed to read message count")
	}

	if numMessages <= 0 {
		return nil, util.NewFatalClientErr(err, "E_BAD_BODY",
			fmt.Sprintf("MPUB invalid message count %d", numMessages))
	}

	messages := make([]*Message, 0, numMessages)
	for i := int32(0); i < numMessages; i++ {
		messageSize, err := readLen(r, tmp)
		if err != nil {
			return nil, util.NewFatalClientErr(err, "E_BAD_MESSAGE",
				fmt.Sprintf("MPUB failed to read message(%d) body size", i))
		}

		if messageSize <= 0 {
			return nil, util.NewFatalClientErr(nil, "E_BAD_MESSAGE",
				fmt.Sprintf("MPUB invalid message(%d) body size %d", i, messageSize))
		}

		if int64(messageSize) > maxMessageSize {
			return nil, util.NewFatalClientErr(nil, "E_BAD_MESSAGE",
				fmt.Sprintf("MPUB message too big %d > %d", messageSize, maxMessageSize))
		}

		msgBody := make([]byte, messageSize)
		_, err = io.ReadFull(r, msgBody)
		if err != nil {
			return nil, util.NewFatalClientErr(err, "E_BAD_MESSAGE", "MPUB failed to read message body")
		}

		messages = append(messages, NewMessage(<-idChan, msgBody))
	}

	return messages, nil
}
