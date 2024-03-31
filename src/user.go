package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	MyMsg  chan string
	myconn net.Conn
	server *Server
	done   chan bool
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := User{
		Name:   userAddr,
		Addr:   userAddr,
		MyMsg:  make(chan string),
		myconn: conn,
		server: server,
		done:   make(chan bool),
	}
	go user.UserListening()
	return &user
}

func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.MapUsers[u.Name] = u
	u.server.mapLock.Unlock()
	u.server.BoradCast(u, "已上线\n")

}

func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.MapUsers, u.Name)
	u.server.mapLock.Unlock()
	u.server.BoradCast(u, "已下线")
	u.done <- true
	close(u.done)

}

func (u *User) Dealmsg(msg string) {
	if msg == "who" {
		u.Checkusers(msg)
	} else if len(msg) > 6 && msg[:6] == "rename" {
		u.Rename(msg)
	} else if len(msg) > 2 && msg[:2] == "to" {
		u.SendTOUser(msg)
	} else {
		u.server.BoradCast(u, msg)
	}

}

func (u *User) Checkusers(msg string) {
	u.server.mapLock.Lock()
	for _, value := range u.server.MapUsers {
		onlineMsg := fmt.Sprintf("[用户 %s ]:在线", string(value.Name))
		u.SendMsg(onlineMsg)
	}
	u.server.mapLock.Unlock()

}

func (u *User) Rename(msg string) {
	oldname := u.Name
	newname := msg[7:]
	u.server.mapLock.Lock()
	if _, ok := u.server.MapUsers[newname]; ok {
		u.SendMsg(fmt.Sprintf("%s用户名已存在", newname))
	} else {
		delete(u.server.MapUsers, u.Name)
		u.Name = newname
		u.server.MapUsers[u.Name] = u
		u.SendMsg(fmt.Sprintf("用户%s 更改名字为 %s", oldname, newname))
	}
	u.server.mapLock.Unlock()
}

func (u *User) SendMsg(msg string) {
	_, err := u.myconn.Write([]byte(msg))
	if err != nil {
		fmt.Println("User发送数据到客户端失败:", err)
		return
	}
}
func (u *User) SendTOUser(msg string) {
	srcname := u.Name
	recvname := strings.Split(msg, " ")[1]
	contexts := strings.Split(msg, " ")[2]
	recvUser, ok := u.server.MapUsers[recvname]
	if !ok {
		u.SendMsg(fmt.Sprintf("用户[%s]不存在", recvname))
		return
	}
	recvUser.SendMsg(fmt.Sprintf("用户[%s]私聊:%s", srcname, string(contexts)))
}

func (u *User) UserListening() {
	for {
		select {
		case msg := <-u.MyMsg:
			u.SendMsg(msg)
		case <-u.done:
			close(u.MyMsg)
			fmt.Printf("用户[%s]下线，结束监听\n", u.Name)
			return
		}
	}

}
