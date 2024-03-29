package main

import (
	"fmt"
	"net"
	"sync"
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
	user := NewUser(conn)
	s.mapLock.Lock()
	s.MapUsers[user.Name] = user
	s.mapLock.Unlock()
	s.BoradCast(user, "已上线")
	select {}
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
			client.C <- msg
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
