//go:build !windows
// +build !windows

package application

import (
	"os/signal"
	"sync"
	"syscall"
)

// SignalNotify 注册监听的信号
func (a *Application) signalNotify() {
	signal.Notify(a.signalChan, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
}

// checkSysSignal 检查信号
func (a *Application) checkSysSignal(wg *sync.WaitGroup) {
	for {
		select {
		case s := <-a.signalChan:
			switch s {
			case syscall.SIGUSR1:
				a.reload()
			default:
				a.exit()
				wg.Done()
			}
		}
	}
}
