package application

import (
	"fly/internal/config"
	"fly/internal/datetime/current"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"time"
)

const TickerTime = time.Millisecond * 100

// IApp 应用程序接口
type IApp interface {
	OnInit() error
	OnReload() error
	OnTick(lastMs, nowMs int64)
	OnExit()
}

type Application struct {
	app IApp

	lastTickTime int64
	signalChan   chan os.Signal // 信号监听通道
}

/*
NewApplication 创建一个应用程序
configFilePath --> 配置文件地址
*/
func NewApplication(handler IApp) (*Application, error) {
	ret := &Application{
		signalChan:   make(chan os.Signal, 1),
		app:          handler,
		lastTickTime: 0,
	}

	// 初始化配置
	var err error
	err = config.Init()
	if err != nil {
		return nil, err
	}

	// 开启Pprof
	if config.Global.App.Pprof {
		ret.PProf(strconv.Itoa(int(config.Global.App.PprofPort)))
	}

	// 开启信号通知
	ret.signalNotify()
	return ret, nil
}

func (a *Application) exit() {
	a.app.OnExit()
}

func (a *Application) reload() error {
	return a.app.OnReload()
}

func (a *Application) tick(lastMs, nowMs int64) {
	a.app.OnTick(lastMs, nowMs)
}

func (a *Application) Run() {
	wg := &sync.WaitGroup{}

	err := a.app.OnInit()
	if err != nil {
		// 初始化失败, 直接 panic()
		panic(fmt.Sprintf("应用初始化失败: %v", err))
	}

	wg.Add(1)
	// 监听 信号
	go a.checkSysSignal(wg)
	// 毫秒级别的定时器
	ticker := time.NewTicker(TickerTime)
	go func() {
		for range ticker.C {
			nowMs := current.UnixMill()
			a.tick(a.lastTickTime, nowMs)
			a.lastTickTime = nowMs
		}
	}()
	wg.Wait()
	os.Exit(0)
}

// PProf 开启监听
func (a *Application) PProf(listenPort string) {
	go func() {
		err := http.ListenAndServe(":"+listenPort, nil)
		if err != nil {
			return
		}
	}()
	return
}
