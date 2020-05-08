package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type clientInfo struct {
	addr    string
	name    string
	sendMsg chan messageInfo
	exit    chan struct{}
	conn    net.Conn
}

type messageInfo struct {
	from    string
	message string
}

// 全局的channel 用于收集客户端发来的消息
var totalMessage chan messageInfo

// 所有连接到的客户端都注册到这里
var client map[string]clientInfo

func init() {
	totalMessage = make(chan messageInfo)
	client = make(map[string]clientInfo)
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Printf("start listen failed, error:%v\n", err)
		return
	}

	// 全局的广播先工作起来
	go BroadcastMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept connection failed, error:%v\n", err)
			continue
		}

		// 客户端连接过来，注册
		addr := conn.RemoteAddr().String()
		log.Printf("client %s connected", addr)
		client[addr] = clientInfo{
			addr:    addr,
			name:    addr,
			sendMsg: make(chan messageInfo),
			exit:    make(chan struct{}),
			conn:    conn,
		}

		// 处理读消息
		go connReadHandle(client[addr])
		// 处理写消息
		go connWriteHandle(client[addr])
	}
}

// 处理客户端退出
func clientExit(c clientInfo) {
	_ = c.conn.Close()
	delete(client, c.addr)
	c.exit <- struct{}{}
}

// 接收客户端消息
func connReadHandle(c clientInfo) {
	defer clientExit(c)

	for {
		buff := make([]byte, 4096)
		n, err := c.conn.Read(buff)
		if err == io.EOF {
			msg := fmt.Sprintf("%s : user %s exit", c.addr, c.name)
			totalMessage <- messageInfo{
				from:    c.addr,
				message: msg,
			}
			return
		}
		if err != nil {
			fmt.Printf("recive message for %s failed, error:%v\n", c.addr, err)
			return
		}
		fmt.Printf("recive from [%s], message:[%s]\n", c.addr, string(buff[:n]))

		// 接收到的消息，发给广播channel
		msg := fmt.Sprintf("[%s]:%s", c.name, string(buff[:n]))
		totalMessage <- messageInfo{
			from:    c.addr,
			message: msg,
		}
	}
}

// 发消息给客户端
func connWriteHandle(c clientInfo) {
	for {
		select {
		case message := <-c.sendMsg:
			if message.from != c.addr {
				_, _ = c.conn.Write([]byte(message.message))
			}
		case <-c.exit:
			log.Printf("close connWriteHandle")
			return
		}
	}
}

// 广播消息给所有的client
func BroadcastMessage() {
	for {
		msg := <-totalMessage

		for _, c := range client {
			c.sendMsg <- msg
		}
	}
}
