package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"time"
)

// createWsClient 创建一个连接
func createWsClient(addr string) (*websocket.Conn, error) {
	var err error
	var conn *websocket.Conn
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	dialer := &websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
	}

	conn, _, err = dialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// createWssClient 创建一个连接
func createWssClient(addr string) (*websocket.Conn, error) {
	var err error
	var conn *websocket.Conn
	u := url.URL{Scheme: "wss", Host: addr, Path: "/ws"}
	rootPath := "/mnt/hgfs/meta/xingqiujueqi-server/ca/"
	baseName := "192.168.1.170"
	//rootPath := "D:\\meta\\xingqiujueqi-server\\ca\\"
	dialer := &websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{RootCAs: loadCA(rootPath + baseName + ".crt")},
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	conn, _, err = dialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func loadCA(caFile string) *x509.CertPool {
	pool := x509.NewCertPool()

	if ca, e := os.ReadFile(caFile); e != nil {
		panic(e)
	} else {
		pool.AppendCertsFromPEM(ca)
	}
	return pool
}

func main() {
	conn, err := createWsClient("192.168.1.170:80")
	if err != nil {
		panic(err)
	}
	err = conn.WriteMessage(websocket.BinaryMessage, []byte("111111"))
	if err != nil {
		panic(err)
	}
	mt, data, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	fmt.Printf("client : %v, %s\n", mt, data)
	err = conn.Close()
	if err != nil {
		panic(err)
	}
}
