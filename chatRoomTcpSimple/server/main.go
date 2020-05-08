package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Printf("start listen failed, error:%v\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept connection failed, error:%v\n", err)
			continue
		}

		go connHandle(conn)
	}
}

func connHandle(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()

	for {
		buff := make([]byte, 4096)
		n, err := conn.Read(buff)
		if err == io.EOF {
			fmt.Printf("client %s closed\n", addr)
			return
		}
		if err != nil {
			fmt.Printf("recive message for %s failed, error:%v\n", addr, err)
			return
		}
		fmt.Printf("recive from [%s], message:[%s]\n", addr, string(buff[:n]))
		conn.Write(buff[:n])
	}
}
