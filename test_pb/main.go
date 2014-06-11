// test_pb project main.go
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	"fmt"
	"log"
	"net"
	"test_pb/myproto"
	"unsafe"
)

func main() {
	fmt.Println("Hello World!")

	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	req := &myproto.RegReq{
		Label: proto.String("hello"),
		Type:  proto.Int32(17)}

	var buf [4]byte
	*((*int)(unsafe.Pointer(&buf))) = proto.Size(req)

	conn.Write(buf[:])
	bytes, err := proto.Marshal(req)
	conn.Write(bytes)

	ch := make(chan int)
	<-ch

}
