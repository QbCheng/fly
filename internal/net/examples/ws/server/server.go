package main

import (
	"fly/internal/net"
	"fmt"
	"time"
)

type wsServer struct{}

func (w *wsServer) OnConn(sess net.IConn) error {
	fmt.Printf("OnConn : %d, %s\n", sess.Id(), sess.RemoteAddr())
	return nil
}

func (w *wsServer) OnPacket(sess net.IConn, data []byte) {
	//fmt.Printf("OnPacket : %d, %s\n", id, string(data))
	fmt.Printf("OnPacket : %d, %s, %s\n", sess.Id(), sess.RemoteAddr(), data)
	// 将消息反射回客户端
	_ = sess.WriteData(&net.WritePacket{
		Uid:    0,
		Cmd:    0,
		Header: []byte(sess.RemoteAddr()),
		Body:   data,
	})
	return
}

func (w *wsServer) OnClose(sess net.IConn) {
	fmt.Printf("OnClose : %d, %s\n", sess.Id(), sess.RemoteAddr())
}

func (w *wsServer) OnSend(sess net.IConn, uid uint64, cmd uint32, sendResult error) {
	if sendResult != nil {
		fmt.Printf("OnSend 发送失败 :%d, %v\n", sess.Id(), sendResult)
	} else {
		fmt.Printf("OnSend 发送成功 : %d, %d, %d\n", sess.Id(), uid, cmd)
	}
}

func WsConfig() net.Options {
	option := net.DefaultWsOptions()
	option.Ip = "0.0.0.0"
	option.Port = "18888"
	return option
}

func main() {
	_, err := net.NewWebSocketSvr(WsConfig(), &wsServer{})
	if err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Minute)
}
