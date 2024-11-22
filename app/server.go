package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	//"golang.org/x/text/cases"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6378")
	if err != nil {
		fmt.Println("Failed to bind to port 6378")
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Implement the handling of the incoming connection
	defer conn.Close()
	cache := make(map[string]string)
	duration := make(map[string]string)

	for {
		buf := make([]byte, 128)
		_, err := conn.Read(buf[:])
		if errors.Is(err, io.EOF) {
			fmt.Println("Client closed the connections:", conn.RemoteAddr())
			break
		}
		if err != nil {
			fmt.Println("Error reading connection: ", err.Error())
			os.Exit(1)
		}
		//	if(res==([]byte("PING"))){
		//conn.Write([]byte("+PONG\r\n"))
		//}
		fmt.Println(strings.Fields(string(buf)))
		v := strings.Fields(string(buf))
		var command []string
		if v[0][:1] == "*" {
			length, _ := strconv.Atoi(v[0][1:])
			for i := 0; i < length; i++ {
				command = append(command, v[i*2+2])
				//fmt.Println(v[i*2+2])
			}

			//fmt.Println(length, command)
			fmt.Println(time.Now())
			switch command[0] {
			case "PING":
				conn.Write([]byte("+PONG\r\n"))
			case "ECHO":
				conn.Write([]byte("+" + command[1] + "\r\n"))
			case "COMMAND":
				conn.Write([]byte("+" + command[1] + "\r\n"))
			case "SET":
				conn.Write([]byte("+OK\r\n"))
				cache[command[1]] = command[2]
				if len(command) > 3 {
					delay, _ := strconv.Atoi(command[4])
					duration[command[1]] = strconv.Itoa(int(time.Now().UnixMilli()) + delay)
				}
			case "GET":
				if val, ok := cache[command[1]]; ok {
					//fmt.Println("$" + string(3) + "\r\n" + val + "\r\n")
					if val2, ok := duration[command[1]]; ok {
						dur, _ := strconv.Atoi(val2)
						if dur > int(time.Now().UnixMilli()) {
							conn.Write([]byte("$" + strconv.Itoa(len(val)) + "\r\n" + val + "\r\n"))
						} else {
							fmt.Println("Key expired")
							fmt.Println(dur, time.Now().UnixMilli())
							conn.Write([]byte("$-1\r\n"))
						}
						//fmt.Println("$" + string(3) + "\r\n" + val + "\r\n")
					} else {
						fmt.Println("No duration found")
						conn.Write([]byte("$" + strconv.Itoa(len(val)) + "\r\n" + val + "\r\n"))
					}
				} else {
					fmt.Println("Key not found")
					conn.Write([]byte("$-1\r\n"))
				}
			}
		}
	}
}

/*func parseresp(message byte[]){
	if message == []byte("PING"){
		return []byte("+PONG\r\n")
	}



}
*/
