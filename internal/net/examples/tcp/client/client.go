package main

import (
	"bufio"
	"fly/internal/net/examples/tcp/packet"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:18888")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		// 读取用户输入
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')

		packet := packet.CSPacket{
			Header: packet.CSPacketHeader{
				Version: 1,
				Uid:     1,
				BodyLen: uint32(len(message)),
			},
			Body: []byte(message),
		}
		packetByte := packet.ToBytes()

		// 发送消息给服务器
		_, err = conn.Write(packetByte)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		buf := make([]byte, 1024)
		_, _ = conn.Read(buf)
		fmt.Printf("Server response: %s\n", buf)
	}
}
