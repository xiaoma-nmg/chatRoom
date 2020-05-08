package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Printf("connect server failed, error:%v\n", err)
		return
	}
	defer conn.Close()

	go func() {
		fmt.Println("Please input something:")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			conn.Write([]byte(scanner.Text()))
		}
	}()

	for {
		buff := make([]byte, 4096)
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("read message from server failed, error:%v\n", err)
		}
		fmt.Println(string(buff[:n]))
	}
}
