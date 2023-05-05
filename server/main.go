package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"net"
)

var clientList = list.New()

func sendToAll(msg string) {
	for elem := clientList.Front(); elem != nil; elem = elem.Next() {
		fmt.Fprint(elem.Value.(net.Conn), msg)
	}
}

func handle(c net.Conn) {
	elem := clientList.PushBack(c)
	defer clientList.Remove(elem)
	reader := bufio.NewReader(c)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		sendToAll(msg)
	}
	log.Println("a client left")
}

func main() {
	ln, err := net.Listen("tcp", ":1922")

	if err != nil {
		panic(err)
	}

	log.Println("Server started at port 1922")
	for {
		conn, err := ln.Accept()

		if err == nil {
			go handle(conn)
		} else {
			log.Println(err)
		}
	}
}
