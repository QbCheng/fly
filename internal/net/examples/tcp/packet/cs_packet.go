package packet

import (
	"encoding/binary"
	"fly/internal/net"
)

var _ net.IPacketInfo = (*CSPacket)(nil)

// CSPacket 服务器和客户端之间的整包
type CSPacket struct {
	Header CSPacketHeader
	Body   []byte
}

// HeaderSize 服务器和客户端之间包头长度
func (h *CSPacket) HeaderSize() int {
	return h.Header.HeaderSize()
}

// BodySizeWithHeader 通过 包头的 []byte 获取当前包 body 长度
func (h *CSPacket) BodySizeWithHeader(header []byte) int {
	return int(binary.BigEndian.Uint32(header[h.HeaderSize()-4:]))
}

func (h *CSPacket) ToBytes() []byte {
	message := make([]byte, h.HeaderSize()+len(h.Body))
	copy(message, h.Header.ToBytes())
	copy(message[h.HeaderSize():], h.Body)
	return message
}

// CSPacketHeader 服务器和客户端之间的包头
type CSPacketHeader struct {
	Version  uint16 // 客户端版本
	PassCode uint16 // 简单秘钥.
	Seq      uint32 // 序列号. 当前没有使用

	Uid uint64 // Uid

	LoginPin uint32 // 登录码. 由网关服分配. 作为该用户本次登录的唯一码
	Cmd      uint32 // 执行的指令.

	BodyLen uint32 // 包长
}

// HeaderSize 服务器和客户端之间包头长度
func (h *CSPacketHeader) HeaderSize() int {
	return 28
}

func (h *CSPacketHeader) From(b []byte) {
	pos := 0
	h.Version = binary.BigEndian.Uint16(b[pos:])
	pos += 2
	h.PassCode = binary.BigEndian.Uint16(b[pos:])
	pos += 2
	h.Seq = binary.BigEndian.Uint32(b[pos:])
	pos += 4
	h.Uid = binary.BigEndian.Uint64(b[pos:])
	pos += 8
	h.LoginPin = binary.BigEndian.Uint32(b[pos:])
	pos += 4
	h.Cmd = binary.BigEndian.Uint32(b[pos:])
	pos += 4
	h.BodyLen = binary.BigEndian.Uint32(b[pos:])
	pos += 4
}

func (h *CSPacketHeader) To(b []byte) {
	pos := uintptr(0)
	binary.BigEndian.PutUint16(b[pos:], h.Version)
	pos += 2
	binary.BigEndian.PutUint16(b[pos:], h.PassCode)
	pos += 2
	binary.BigEndian.PutUint32(b[pos:], h.Seq)
	pos += 4
	binary.BigEndian.PutUint64(b[pos:], h.Uid)
	pos += 8
	binary.BigEndian.PutUint32(b[pos:], h.LoginPin)
	pos += 4
	binary.BigEndian.PutUint32(b[pos:], h.Cmd)
	pos += 4
	binary.BigEndian.PutUint32(b[pos:], h.BodyLen)
	pos += 4
}

func (h *CSPacketHeader) ToBytes() []byte {
	bytes := make([]byte, h.HeaderSize())
	h.To(bytes)
	return bytes
}
