package net

import (
	"time"
)

// INet 网络接口
type INet interface {
	Run() error
	WriteData(id uint64, uid uint64, cmd uint32, header []byte, body []byte) error
	SafeCloseWithConnId(id uint64) error
}

// IConnEvent 连接事件
type IConnEvent interface {
	OnConn(session IConn) error
	OnPacket(session IConn, data []byte)
	OnClose(session IConn)
	OnSend(session IConn, uid uint64, cmd uint32, err error)
}

// IConn Conn 统一接口
type IConn interface {
	Id() uint64
	LocalAddr() string
	RemoteAddr() string // 如果使用了代理服务器, 该接口将返回 `客户端IP地址 + 代理服务器和后端服务器实际连接的端口`
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
	WriteData(*WritePacket) error
}

// WritePacket 写入的数据包
type WritePacket struct {
	Header []byte
	Body   []byte
	Cmd    uint32
	Uid    uint64
}
