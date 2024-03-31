package main

import (
	"fmt"
	"io"
	"net"
	"strings"
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

	timeout := 60 * time.Second
	readCh := make(chan string)
	errCh := make(chan error)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 || err != nil {
				if n == 0 {
					err = io.EOF
				}
				errCh <- err
				return
			}
			msg := strings.TrimSpace(string(buf[:n]))
			readCh <- msg

		}
	}()

	for {
		select {
		case msg := <-readCh:
			fmt.Printf("用户[%s]处理消息\n", user.Name)
			user.Dealmsg(string(msg))
		case err := <-errCh:
			if err == io.EOF {
				fmt.Printf("用户[%s]断开连接\n", user.Name)
			} else {
				fmt.Printf("端口读消息出错: %s\n", err)
			}
			user.Offline()
			return
		case <-time.After(timeout):
			fmt.Printf("用户[%s]长时间未响应\n", user.Name)
			user.SendMsg("长时间未响应，已强制下线")
			user.Offline()
			return
		}
	}

}

func (s *Server) BoradCast(user *User, msg string) {
	boradMsg := fmt.Sprintf("[用户 %s]:%s", user.Name, msg)
	s.Msg <- boradMsg
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
