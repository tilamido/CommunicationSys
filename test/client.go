package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()
	fmt.Print("连接成功\n")
	go handleStdin(ctx, conn, cancel)

	handleServerMessages(ctx, conn)
}

func handleStdin(ctx context.Context, conn net.Conn, cancel context.CancelFunc) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !scanner.Scan() {
				cancel()
				return
			}
			msg := scanner.Text()
			if msg == "exit" {
				cancel()
				return
			}
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("Error writing to server:", err)
				cancel()
				return
			}
		}
	}
}

func handleServerMessages(ctx context.Context, conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("客户端读取数据失败:", err)
			} else {
				fmt.Println("服务端主动断开连接")
			}
			os.Exit(0) // 直接退出程序
		}
		message := string(buffer[:n])
		fmt.Printf("Received: %s\n", message)
	}
}
