// 文件名：client/main.go
package main

import (
	"fmt"
	"net"
)

func main() {

	// 建立 TCP 连接
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	for {
		// 读取来自服务器的响应
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from server:", err.Error())
			break
		}
		// 打印来自服务器的消息
		fmt.Printf("Received: %s\n", string(buffer[:n]))
	}

}
