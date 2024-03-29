package main

import (
	"fmt"
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	myconn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		myconn: conn,
	}

	go user.UserListening()
	return &user
}

func (u *User) UserListening() {

	for {
		msg := <-u.C
		fmt.Println("UserListening\n", "msg: ", msg)
		// _, err := u.myconn.Write([]byte(msg + "\n"))
		// if err != nil {
		// 	fmt.Println("Error writing:", err.Error())
		// 	return
		// }
	}

}
