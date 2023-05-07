package net

import (
	"time"
)

type Options struct {
	Ip                  string
	Port                string
	SessionReadTimeout  time.Duration // 读取超时
	SessionWriteTimeout time.Duration // 写入超时
	Pattern             string        // http pattern method
	CertFile, KeyFile   string        // TLS
	Logger              Logger
	AllLogger           bool // AllLogger 开启所有日志
}

func DefaultOptions() Options {
	return Options{
		Ip:                  "127.0.0.1",
		Port:                "8888",
		SessionReadTimeout:  30 * time.Second,
		SessionWriteTimeout: 30 * time.Second,
		CertFile:            "",
		KeyFile:             "",
		Pattern:             "/ws",
		Logger:              NewSimpleLog(),
		AllLogger:           true,
	}
}
