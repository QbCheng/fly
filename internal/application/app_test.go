package application

import (
	"fmt"
	"testing"
)

var _ IApp = (*ExampleSvrImpl)(nil)

type ExampleSvrImpl struct {
}

func (a *ExampleSvrImpl) OnInit() error {
	fmt.Println("ExampleSvrImpl : OnInit")
	return nil
}

func (a *ExampleSvrImpl) OnExit() {
	fmt.Println("ExampleSvrImpl : OnExit")
}

func (a *ExampleSvrImpl) OnReload() error {
	fmt.Println("ExampleSvrImpl : OnReload")
	return nil
}

func (a *ExampleSvrImpl) OnTick(lastMs, nowMs int64) {
	fmt.Printf("ExampleSvrImpl : tick. {%d, %d}\n", lastMs, nowMs)
}

func TestNewApplication(t *testing.T) {
	app, err := NewApplication(&ExampleSvrImpl{})
	if err != nil {
		panic(err)
	}
	// 开启 pprof
	app.PProf("18888")
	fmt.Println("启动成功")
	app.Run()
}
