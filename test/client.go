package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on :8888")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// 启动一个goroutine来处理连接
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	// 向客户端发送消息
	_, err := conn.Write([]byte("收到请求\n"))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		return
	}
	// 可以继续读取客户端发送的数据，或者做其他处理
	// ...
}
