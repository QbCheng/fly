//go:build windows
// +build windows

package application

import (
	"os/signal"
	"sync"
	"syscall"
)

// signalNotify 注册监听的信号
func (a *Application) signalNotify() {
	signal.Notify(a.signalChan, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
}

func (a *Application) checkSysSignal(wg *sync.WaitGroup) {
	for {
		select {
		case s := <-a.signalChan:
			switch s {
			default:
				a.exit()
				wg.Done()
			}
		}
	}
}
