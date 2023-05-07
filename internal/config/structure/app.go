package structure

import (
	"fly/internal/config/clientVersion"
)

type App struct {
	AppName string `  json:"app_name" yaml:"app_name" mapstructure:"app_name"` // app 名字
	AppId   int    `json:"app_id" yaml:"app_id" mapstructure:"app_id"`         // app id

	ServerId          string `json:"server_id" yaml:"server_id" mapstructure:"server_id"` // 服务ID  busId 对象的字符串格式
	ServerIdUint64    uint64 `json:"-" yaml:"-" mapstructure:"-"`                         // serverIdUint64 是 ServerId Uint64 版本. 通过算法解析得到
	SeverType         uint64 `json:"-" yaml:"-" mapstructure:"-"`                         // severType 是 ServerId 解析之后的服务类型
	SeverTypeInstance uint64 `json:"-" yaml:"-" mapstructure:"-"`                         // severTypeInstance 是 ServerId 解析之后的服务类型的实例ID

	ClientVersion       string `json:"client_version" yaml:"client_version" mapstructure:"client_version"` // 客户端版本号. 用于客户端版本校验, 客户端版本号 和 服务器保存的客户端版本不同时, 不允许进入游戏.
	ClientVersionUint16 uint16 `json:"-" yaml:"-" mapstructure:"-"`                                        // 客户端版本号的 uint16 版本. 通过算法解析得到

	Development int `json:"development" yaml:"development" mapstructure:"development"` // 开发环境 1: 开发环境. 0:正式环境

	Pprof     bool   `json:"pprof" yaml:"pprof" mapstructure:"pprof"` // pprof 监听的端口
	PprofPort uint64 `json:"-" yaml:"-" mapstructure:"-""`            // pprof 监听的端口
}

func (a *App) Pretreatment() error {
	var err error
	// 转换客户端版本
	a.ClientVersionUint16, err = clientVersion.FormatClientVersion(a.ClientVersion)
	if err != nil {
		return err
	}

	// 服务Id 的busId类型
	bid, err := busId.NewBusId(a.ServerId)
	if err != nil {
		return err
	}

	a.ServerIdUint64 = bid.Iid()
	a.SeverType = bid.Get(busId.ParseResultTyp)
	a.SeverTypeInstance = bid.Get(busId.ParseResultInstance)

	// pprof 是否开启
	if a.Pprof {
		a.PprofPort = 10000 + a.SeverType*100 + a.SeverTypeInstance
	}
	return nil
}
