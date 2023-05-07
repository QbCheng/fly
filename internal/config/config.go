package config

import "fly/internal/config/structure"

type Wrap struct {
	App                  structure.App                  `mapstructure:"app" json:"app" yaml:"app"`
	Redis                []structure.Redis              `mapstructure:"redis" json:"redis" yaml:"redis"`
	BuriedPointCollector structure.BuriedPointCollector `mapstructure:"buried_point_collector" json:"buried_point_collector" yaml:"buried_point_collector"`
	Mysql                []structure.Mysql              `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	LogCollector         structure.LogCollector         `mapstructure:"log_collector" json:"log_collector" yaml:"log_collector"`
	Router               structure.Router               `mapstructure:"router" json:"router" yaml:"router"`

	HttpServer structure.HttpServer `mapstructure:"http_server" json:"http_server" yaml:"http_server"`
	MainServer structure.MainServer `mapstructure:"main_server" json:"main_server" yaml:"main_server"`
	ConnServer structure.ConnServer `mapstructure:"conn_server" json:"conn_server" yaml:"conn_server"`
	GateServer structure.GateServer `mapstructure:"gate_server" json:"gate_server" yaml:"gate_server"`
}

// Pretreatment 对配置进行预处理
func (w *Wrap) Pretreatment() error {
	if err := w.App.Pretreatment(); err != nil {
		return err
	}

	for i := range w.Redis {
		if err := w.Redis[i].Pretreatment(); err != nil {
			return err
		}
	}

	for i := range w.Mysql {
		if err := w.Mysql[i].Pretreatment(); err != nil {
			return err
		}
	}
	return nil
}
