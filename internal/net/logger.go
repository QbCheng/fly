package net

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Printf(string, ...interface{})
}

type SimpleLog struct {
	log *log.Logger
}

func NewSimpleLog() *SimpleLog {
	return &SimpleLog{
		log: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

func (d *SimpleLog) Printf(layout string, args ...interface{}) {
	_ = d.log.Output(3, fmt.Sprintf(layout, args...))
}
