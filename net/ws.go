package net

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strings"
	"sync"
	"time"
)

var _ IConn = (*wsConn)(nil)

type wsConn struct {
	id   uint64          // Session 唯一ID
	conn *websocket.Conn // websocket 连接

	realClientIp string // 真实客户端ip
	port         string // 端口
	remoteAddr   string // 真实客户端ip + port

	chWrite chan *WritePacket
	close   *atomic.Bool // 关闭状态
	logger  Logger
	cb      IConnEvent
}

// NewWsConn 创建
func NewWsConn(id uint64, conn *websocket.Conn, log Logger) *wsConn {
	return &wsConn{
		id:      id,
		conn:    conn,
		chWrite: make(chan *WritePacket, 1),
		close:   atomic.NewBool(false),
		logger:  log,
	}
}

// Id 唯一标识
func (s *wsConn) Id() uint64 {
	return s.id
}

// LocalAddr 本地地址
func (s *wsConn) LocalAddr() string {
	return s.conn.LocalAddr().String()
}

// RemoteAddr 远端地址
func (s *wsConn) RemoteAddr() string {
	return s.remoteAddr
}

// SetReadDeadline 设置 底层连接 读取一行数据的数据超时时间
func (s *wsConn) SetReadDeadline(DeadTime time.Time) error {
	return s.conn.SetReadDeadline(DeadTime)
}

// SetWriteDeadline 设置 底层连接 写入一行数据的数据超时时间
func (s *wsConn) SetWriteDeadline(DeadTime time.Time) error {
	return s.conn.SetWriteDeadline(DeadTime)
}

// SetCloseState 将 当前连接 设置为关闭状态
func (s *wsConn) SetCloseState() {
	s.close.Store(true)
}

// IsClose 当前连接 是关闭的
func (s *wsConn) IsClose() bool {
	return s.close.CAS(true, true)
}

// WriteData 写入一个消息
func (s *wsConn) WriteData(packet *WritePacket) error {
	select {
	case s.chWrite <- packet:
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout in 10 seconds")
	}
	return nil
}

func (s *wsConn) tag() string {
	return fmt.Sprintf("ws->%d, remote->%s |", s.id, s.RemoteAddr())
}

func (s *wsConn) Error(layout string, args ...interface{}) {
	s.logger.Printf(s.tag()+"error | "+layout, args...)
}

type WebsocketSvr struct {
	config Options
	cb     IConnEvent

	upgrade *websocket.Upgrader
	server  *http.Server

	idLock   sync.RWMutex
	id       uint64
	idToConn map[uint64]*wsConn
}

func NewWebSocketSvr(config Options, handler IConnEvent) (INet, error) {
	s := &WebsocketSvr{
		config: config,
		// Upgrader 指定HTTP连接升级到WebSocket连接.
		// Upgrader 上的方法是并发安全的
		upgrade: &websocket.Upgrader{
			HandshakeTimeout: 1 * time.Minute,
			ReadBufferSize:   1024 * 10,
			WriteBufferSize:  1024 * 10,
			// 解决跨域问题
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		cb: handler,
		server: &http.Server{
			Addr:    config.Ip + ":" + config.Port,
			Handler: nil,
		},

		idLock:   sync.RWMutex{},
		id:       0,
		idToConn: map[uint64]*wsConn{},
	}
	err := s.Run()
	return s, err
}

// setRemoteAddr 设置远端地址.
// 如果使用了代理, 远端地址 = 真实客户端IP + 代理客户端和后台服务器创建的连接
// 没有使用代理, 则是使用的真实的远端地址
func (s *WebsocketSvr) setRemoteAddr(localConn *wsConn, r *http.Request) {
	// 判断是否使用了代理
	if r.Header.Get("X-Forwarded-For") != "" {
		// 获取真实客户端IP地址
		localConn.realClientIp = r.Header.Get("X-Real-IP")
		if strs := strings.Split(r.RemoteAddr, ":"); len(strs) >= 2 {
			localConn.port = strs[1]
		}
		// 实际上, 后端服务器是与代理服务器相连接, 并非与客户机相连接
		localConn.remoteAddr = localConn.realClientIp + ":" + localConn.port
	} else {
		localConn.remoteAddr = r.RemoteAddr
		if strs := strings.Split(r.RemoteAddr, ":"); len(strs) >= 2 {
			localConn.realClientIp = strs[0]
			localConn.port = strs[1]
		}
	}
	return
}

// Run 启动 Websocket server
func (s *WebsocketSvr) Run() error {
	http.HandleFunc(s.config.Pattern, func(w http.ResponseWriter, r *http.Request) {
		// 将 HTTP 连接 提升 到 Websocket
		nativeConn, err := s.upgrade.Upgrade(w, r, nil)
		if err != nil {
			s.error("Upgrade failed to Upgrade HTTP server connections to the WebSocket protocol. {err : %s}", err.Error())
			return
		}

		s.idLock.Lock()
		defer s.idLock.Unlock()
		s.id++
		localConn := NewWsConn(s.id, nativeConn, s.config.Logger)
		s.setRemoteAddr(localConn, r)
		if err = s.cb.OnConn(localConn); err != nil {
			// 回调失败, 直接关闭该连接.
			if closeErr := nativeConn.Close(); closeErr != nil {
				s.error("OnConn callback failed : %s, Close the current connect failed : %s", err.Error(), closeErr.Error())
			} else {
				s.error("OnConn callback failed : %s, Close the current connect success.", err.Error())
			}
			return
		}
		s.idToConn[s.id] = localConn
		go s.coRead(s.idToConn[s.id])
		go s.coWrite(s.idToConn[s.id])
	})

	// 监听
	go func() {
		if s.config.CertFile == "" || s.config.KeyFile == "" {
			// ws
			s.info("WebSocket Listening on : ws://%s/ws", s.server.Addr)
			err := s.server.ListenAndServe()
			if err != nil {
				s.error("WebSocket Failed to listen. {err : %s}", err)
				return
			}
		} else {
			// wss
			s.info("WebSocket Listening on : wss://%s/ws", s.server.Addr)
			err := s.server.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
			if err != nil {
				s.error("WebSocket Failed to listen. {err : %s}", err)
				return
			}
		}
	}()
	return nil
}

func (s *WebsocketSvr) info(layout string, args ...interface{}) {
	if s.config.AllLogger {
		s.config.Logger.Printf("websocket | info | "+layout, args...)
	}
}

func (s *WebsocketSvr) error(layout string, args ...interface{}) {
	s.config.Logger.Printf("websocket | err | "+layout, args...)
}

// WriteData 写入数据到对应的 id.
func (s *WebsocketSvr) WriteData(id uint64, uid uint64, cmd uint32, header []byte, body []byte) error {
	s.idLock.RLock()
	defer s.idLock.RUnlock()

	sess, exists := s.idToConn[id]
	if !exists {
		return fmt.Errorf("conn doesn't exist")
	}

	if sess.IsClose() {
		return fmt.Errorf("conn : %d is closing", sess.id)
	}

	return sess.WriteData(&WritePacket{
		Cmd:    cmd,
		Uid:    uid,
		Header: header,
		Body:   body,
	})
}

// SafeCloseWithConnId 安全的关闭指定的 id 的连接.
func (s *WebsocketSvr) SafeCloseWithConnId(id uint64) error {
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

// SafeClose 安全的关闭所有的conn.
func (s *WebsocketSvr) SafeClose() {
	s.idLock.Lock()
	defer s.idLock.Unlock()
	for _, sess := range s.idToConn {
		sess.SetCloseState()
		err := sess.WriteData(nil)
		if err != nil {
			s.error("Failed to close the conn : %d. {err :  %v} ", sess.id, err)
		}
	}
	return
}

// SafeRunCoroutine
// 安全的启动一个协程.
// 接收协程的异常, 避免一个协程导致崩溃导致整个服务器崩溃.
func (s *WebsocketSvr) SafeRunCoroutine() {
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

// CoRead 负责读取的协程
func (s *WebsocketSvr) coRead(sess *wsConn) {
	// 捕获异常
	defer s.SafeRunCoroutine()
	defer func() {
		s.cb.OnClose(sess)
		s.clear(sess)
	}()

	for {
		// 设置读取的超时时间
		err := sess.SetReadDeadline(time.Now().Add(s.config.SessionReadTimeout))
		if err != nil {
			break
		}
		// 阻塞读取
		messageType, data, err := sess.conn.ReadMessage()
		if err == nil {
			// 当前只接收二进制数据
			if messageType != websocket.BinaryMessage {
				s.error("WebSocket Connection closed. Websocket reads messages of a type other than BinaryMessage {remoteAddr : %s, messageType : %d} ", sess.RemoteAddr(), messageType)
				break
			}
			// websocket 读取直接是一个完整的包.
			s.cb.OnPacket(sess, data)
		} else {
			websocket.IsCloseError(err, websocket.CloseNormalClosure)
			//s.info("WebSocket Connection closed. {err : %s} ", err)
			if websocket.IsCloseError(
				err,
				// 正常关闭.
				websocket.CloseNormalClosure,
				// 用户使用APP登录. 直接关闭掉游戏.
				websocket.CloseAbnormalClosure,
				// 用户直接点击浏览器关闭按钮
				websocket.CloseGoingAway,
				// app端 关闭时, 会能够该错误
				websocket.CloseNoStatusReceived,
			) {
				break
			} else {
				s.error("WebSocket Connection closed. {err : %s} ", err)
			}
			break
		}
	}
}

// clear 清理数据
func (s *WebsocketSvr) clear(sess *wsConn) {
	s.idLock.Lock()
	defer s.idLock.Unlock()
	close(sess.chWrite)
	delete(s.idToConn, sess.id)
}

// CoWrite 负责写入消息的协程
func (s *WebsocketSvr) coWrite(sess *wsConn) {
	// 捕获异常
	defer s.SafeRunCoroutine()
	for {
		packet, ok := <-sess.chWrite
		if !ok { // chan is closed
			s.info("chanWrite is closed. { conn : %d }", sess.Id())
			break
		}

		// writeData == nil -> 主动关闭
		if packet == nil {
			s.info("A 'nil' is passed to chanWrite to close net. { conn : %d }", sess.Id())
			break
		}

		// 设置 写操作的超时时间
		err := sess.SetWriteDeadline(time.Now().Add(s.config.SessionWriteTimeout))
		if err != nil {
			s.error("Set write deadline error. {err : %s}", err)
			continue
		}

		message := make([]byte, len(packet.Header)+len(packet.Body))
		copy(message, packet.Header)
		copy(message[len(packet.Header):], packet.Body)
		err = sess.conn.WriteMessage(websocket.BinaryMessage, message)
		if err != nil {
			continue
		}
	}
}
