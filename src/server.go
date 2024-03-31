package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip       string
	Port     int
	MapUsers map[string]*User
	mapLock  sync.Mutex
	Msg      chan string
}

// 创建server
func NewSever(ip string, port int) *Server {
	sever := &Server{
		Ip:       ip,
		Port:     port,
		MapUsers: make(map[string]*User),
		Msg:      make(chan string),
	}
	return sever
}

func (s *Server) Handler(conn net.Conn) {
	defer conn.Close()

	user := NewUser(conn, s)
	user.Online()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil {
				if err != io.EOF {
					fmt.Printf("读取错误：%s\n", err)
				} else {
					fmt.Printf("客户端关闭连接：%s\n", conn.RemoteAddr().String())
				}
				return
			}
			msg := buf[:n-1]
			user.Dealmsg(string(msg))
			isLive <- true
		}

	}()

	for {
		select {
		case <-isLive:
			fmt.Printf("用户[%s]处理消息\n", user.Name)
		case <-time.After(20 * time.Second):
			fmt.Printf("用户[%s]长时间未响应\n", user.Name)
			user.SendMsg("长时间未响应，已强制下线")
			user.Offline()
			return
		}
	}

}

func (s *Server) BoradCast(user *User, msg string) {

	sendMsg := "用户-" + user.Name + ":" + msg
	s.Msg <- sendMsg
}
func (s *Server) ListenMsg() {
	for {
		msg := <-s.Msg
		s.mapLock.Lock()
		for _, client := range s.MapUsers {
			client.MyMsg <- msg
		}
		s.mapLock.Unlock()
	}
}
func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Print("net.listen err", err)
		return
	}
	defer listener.Close()
	fmt.Print("服务器启动，等待连接...\n")
	go s.ListenMsg()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("conn err:", err)
			continue
		}
		go s.Handler(conn)
	}

}
