package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"errors"
	"net"
	"io"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6377")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	for{
	buf:=make([]byte,128)
	_, err = conn.Read(buf[:])
	if errors.Is(err, io.EOF) {
	fmt.Println("Client closed the connections:", conn.RemoteAddr())
	break}
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}
//	if(res==([]byte("PING"))){
	conn.Write([]byte("+PONG\r\n"))
//}
	fmt.Println(string(buf))}
}
