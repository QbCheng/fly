package net

import (
	"bytes"
	"fmt"
	"go.uber.org/atomic"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

const (
	kReadBufSize = 1024 * 1024 // 1MB
)

type IPacketInfo interface {
	HeaderSize() int               // 获得包头长度
	BodySizeWithHeader([]byte) int // 通过包头, 获取 包体 长度
}

var _ IConn = (*tcpConn)(nil)

type tcpConn struct {
	id      uint64
	conn    net.Conn
	chWrite chan *WritePacket

	isClose *atomic.Bool
}

func NewTcpConn(id uint64, conn net.Conn) *tcpConn {
	return &tcpConn{
		id:      id,
		conn:    conn,
		chWrite: make(chan *WritePacket, 1),
		isClose: atomic.NewBool(false),
	}
}

func (s *tcpConn) Id() uint64 {
	return s.id
}

func (s *tcpConn) LocalAddr() string {
	return s.conn.LocalAddr().String()
}

func (s *tcpConn) RemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *tcpConn) SetReadDeadline(DeadTime time.Time) error {
	return s.conn.SetReadDeadline(DeadTime)
}
func (s *tcpConn) SetWriteDeadline(DeadTime time.Time) error {
	return s.conn.SetWriteDeadline(DeadTime)
}

func (s *tcpConn) SetCloseState() {
	s.isClose.Store(true)
}

func (s *tcpConn) IsClose() bool {
	return s.isClose.CAS(true, true)
}

// WriteData 写入一个消息
func (s *tcpConn) WriteData(packet *WritePacket) error {
	select {
	case s.chWrite <- packet:
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout in 10 seconds")
	}
	return nil
}

type TcpSvr struct {
	config Options
	cb     IConnEvent

	listener   net.Listener // 监听
	packetInfo IPacketInfo  // Tcp包

	idLock   sync.RWMutex
	id       uint64
	idToConn map[uint64]*tcpConn
}

func NewTcpSvr(config Options, handler IConnEvent, packetInfo IPacketInfo) (INet, error) {
	s := &TcpSvr{
		config:     config,
		cb:         handler,
		idLock:     sync.RWMutex{},
		id:         0,
		idToConn:   map[uint64]*tcpConn{},
		packetInfo: packetInfo,
	}
	var err error
	s.listener, err = net.Listen("tcp", config.Ip+":"+config.Port)
	if err != nil {
		return nil, err
	}

	err = s.Run()
	if err != nil {
		return nil, err
	}
	return s, err
}

func (s *TcpSvr) info(layout string, args ...interface{}) {
	if s.config.AllLogger {
		s.config.Logger.Printf("tcp | info | "+layout, args...)
	}
}

func (s *TcpSvr) error(layout string, args ...interface{}) {
	s.config.Logger.Printf("tcp | err | "+layout, args...)
}

// Run 启动
func (s *TcpSvr) Run() error {
	go func() {
		for {
			// 等待 并返回一个新的连接
			conn, err := s.listener.Accept()
			if err != nil {
				s.error("TCP Accept failure. {err : %s}", err)
				return
			}
			s.onAccept(conn)
		}
	}()
	return nil
}

func (s *TcpSvr) onAccept(conn net.Conn) {
	s.idLock.Lock()
	defer s.idLock.Unlock()
	s.id++
	session := NewTcpConn(s.id, conn)
	if err := s.cb.OnConn(session); err != nil {
		// 回调失败, 直接关闭该连接.
		if closeErr := conn.Close(); closeErr != nil {
			s.error("OnConn callback failed : %s, Close the current connect failed : %s", err.Error(), closeErr.Error())
		} else {
			s.error("OnConn callback failed : %s, Close the current connect success.", err.Error())
		}
		return
	}
	s.idToConn[s.id] = session
	go s.coRead(s.idToConn[s.id])
	go s.coWrite(s.idToConn[s.id])
}

// clear 清理数据
func (s *TcpSvr) clear(sess *tcpConn) {
	s.idLock.Lock()
	defer s.idLock.Unlock()

	// 防止重复关闭
	close(sess.chWrite)
	delete(s.idToConn, sess.id)
}

func (s *TcpSvr) SafeRunCoroutine() {
	if p := recover(); p != nil {
		var printStackTrace = func(err interface{}) string {
			buf := ""
			buf += fmt.Sprintf("stackTrace : {{ %v\n", err)
			for i := 1; ; i++ {
				pc, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				buf += fmt.Sprintf("%s:%d (0x%x)\n", file, line, pc)
			}
			buf += fmt.Sprintf(buf, " }}")
			return buf
		}
		// 打印异常堆栈日志.
		s.error(printStackTrace(p))
	}
}

func (s *TcpSvr) coRead(sess *tcpConn) {
	// 捕获异常
	defer s.SafeRunCoroutine()
	defer func() {
		s.cb.OnClose(sess)
		s.clear(sess)
	}()

	var buff bytes.Buffer
	readBuf := make([]byte, kReadBufSize)
	for {
		// 设置读取的超时时间
		err := sess.SetReadDeadline(time.Now().Add(s.config.SessionReadTimeout))
		if err != nil {
			break
		}
		// 阻塞读取
		readLen, err := sess.conn.Read(readBuf)
		if err == nil {
			buff.Write(readBuf[0:readLen])
			for {
				if buff.Len() >= s.packetInfo.HeaderSize() { // 接收到包头
					bufferData := buff.Bytes()
					bodyLen := s.packetInfo.BodySizeWithHeader(bufferData[:s.packetInfo.HeaderSize()])
					packetLen := s.packetInfo.HeaderSize() + bodyLen
					if len(bufferData) >= packetLen { // 接收到完整的包
						s.cb.OnPacket(sess, bufferData[:packetLen])
						buff.Next(packetLen)
					} else {
						break
					}
				}
				break
			}
		} else {
			if netErr, ok := err.(*net.OpError); ok {
				if netErr.Err.Error() == "use of closed network connection" {
					s.info("An TCP connection closed. {closeAddr : %s, err : %s}", netErr.Addr.String(), netErr.Err)
					break
				}
			}
			if err == io.EOF {
				s.info("An error occurred while reading the TCP connection. {err : io.EOF} ")
				break
			} else {
				s.error("An error occurred while reading the TCP connection. {err : %s} ", err.Error())
				break
			}
		}
	}
}

func WriteAllData(conn net.Conn, packet *WritePacket) error {

	packetLen := len(packet.Header) + len(packet.Body)
	message := make([]byte, packetLen)
	copy(message, packet.Header)
	copy(message[len(packet.Header):], packet.Body)

	bytesWritten := 0

	for bytesWritten < packetLen {
		n, err := conn.Write(message[bytesWritten:])
		if err != nil {
			return err
		}
		bytesWritten += n
	}

	return nil
}

func (s *TcpSvr) coWrite(sess *tcpConn) {
	// 捕获异常
	defer s.SafeRunCoroutine()
	for {
		packet, ok := <-sess.chWrite
		if !ok { // chan is closed
			s.info("chanWrite is closed. { session : %d }", sess.Id())
			break
		}

		// writeData == nil -> 主动关闭
		if packet == nil {
			s.info("A 'nil' is passed to chanWrite to close net. { session : %d }", sess.Id())
			break
		}

		// 设置 写操作的超时时间
		err := sess.SetWriteDeadline(time.Now().Add(s.config.SessionWriteTimeout))
		if err != nil {
			s.error("Set write deadline error. {err : %s}", err)
			s.cb.OnSend(sess, packet.Uid, packet.Cmd, err)
			break
		}

		err = WriteAllData(sess.conn, packet)
		if err != nil {
			s.cb.OnSend(sess, packet.Uid, packet.Cmd, err)
			break
		}

		// 成功发送
		s.cb.OnSend(sess, packet.Uid, packet.Cmd, nil)
	}
}

func (s *TcpSvr) SafeClose() {
	s.idLock.Lock()
	defer s.idLock.Unlock()
	for _, sess := range s.idToConn {
		sess.SetCloseState()
		err := sess.WriteData(nil)
		if err != nil {
			s.error("Failed to close the session : %d. {err :  %v} ", sess.id, err)
		}
	}
	return
}

func (s *TcpSvr) WriteData(id uint64, uid uint64, cmd uint32, header []byte, body []byte) error {
	s.idLock.RLock()
	defer s.idLock.RUnlock()

	sess, exists := s.idToConn[id]
	if !exists {
		return fmt.Errorf("session doesn't exist")
	}

	if sess.IsClose() {
		return fmt.Errorf("session : %d is closed", sess.id)
	}

	return sess.WriteData(&WritePacket{
		Cmd:    cmd,
		Uid:    uid,
		Header: header,
		Body:   body,
	})
}

func (s *TcpSvr) SafeCloseWithConnId(id uint64) error {
	s.idLock.RLock()
	defer s.idLock.RUnlock()
	sess, exists := s.idToConn[id]
	if !exists {
		return fmt.Errorf("connection doesn't exist")
	}

	// 设置关闭状态.
	sess.SetCloseState()
	return sess.WriteData(nil)
}
