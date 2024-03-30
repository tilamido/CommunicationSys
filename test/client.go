package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()
	var wg sync.WaitGroup

	stopCh := make(chan struct{})
	wg.Add(2)

	go func() {
		defer wg.Done()
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("请输入消息：")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from stdin:", err)
				close(stopCh)
				return
			}
			input = strings.TrimSpace(input)
			if input == "exit" {
				close(stopCh)
				return
			}
			if _, err := conn.Write([]byte(input + "\n")); err != nil {
				fmt.Println("Error writing to server:", err)
				close(stopCh)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		for {
			select {
			case <-stopCh:
				return
			default:
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Error reading from server:", err.Error())
					close(stopCh)
					return
				}
				fmt.Printf("Received: %s\n", string(buffer[:n]))
			}
		}
	}()

	<-stopCh
	conn.Close() // 关闭网络连接以确保conn.Read退出阻塞
	wg.Wait()
	fmt.Println("Disconnected from server")
}
