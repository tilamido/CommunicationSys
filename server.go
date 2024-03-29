package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建server
func NewSever(ip string, port int) *Server {
	sever := &Server{
		Ip:   ip,
		Port: port,
	}
	return sever
}
func (this *Server) handler(conn net.Conn) {
	fmt.Print(this.Ip, ":连接成功")
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Print("net.listen err", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("conn err:", err)
			continue
		}

		go this.handler(conn)

	}

}
