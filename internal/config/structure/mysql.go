package structure

import "time"

// Redis 配置项

// Mysql 配置项
type Mysql struct {
	Name                string        `json:"name" yaml:"name" mapstructure:"name"`                            // 实例名
	Addr                string        `json:"addr" yaml:"addr" mapstructure:"addr"`                            // Mysql地址
	User                string        `json:"user" yaml:"user" mapstructure:"user"`                            // 用户名
	Passwd              string        `json:"passwd" yaml:"passwd" mapstructure:"passwd"`                      // 密码
	DBName              string        `json:"db_name" yaml:"db_name" mapstructure:"db_name"`                   // 数据库名
	MaxIdleConn         int           `json:"max_idle_conn" yaml:"max_idle_conn" mapstructure:"max_idle_conn"` // 最大空闲连接
	ConnMaxIdle         string        `json:"conn_max_idle" yaml:"conn_max_idle" mapstructure:"conn_max_idle"` // 3m10s 3分钟10秒
	ConnMaxIdleDuration time.Duration `json:"-" yaml:"-" mapstructure:"-"`                                     //
	MaxOpenConn         int           `json:"max_open_conn" yaml:"max_open_conn" mapstructure:"max_open_conn"` // 最大打开连接数量
}

func (a *Mysql) Pretreatment() error {
	var err error
	a.ConnMaxIdleDuration, err = time.ParseDuration(a.ConnMaxIdle)
	if err != nil {
		return err
	}
	return nil
}
