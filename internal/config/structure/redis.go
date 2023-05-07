package structure

import (
	"time"
)

// Redis 配置项
type Redis struct {
	Name     string `json:"name" yaml:"name" mapstructure:"name"`             // 实例名
	Addr     string `json:"addr" yaml:"addr" mapstructure:"addr"`             // Redis地址 ip:端口, 集群使用 | 进行分割
	Password string `json:"password" yaml:"password" mapstructure:"password"` // Redis账号
	Db       int    `json:"db" yaml:"db" mapstructure:"db"`                   // Redis库

	PoolSize            int           `json:"pool_size" yaml:"pool_size" mapstructure:"pool_size"`                         // Redis连接池大小
	MinIdleConnects     int           `json:"min_idle_connects" yaml:"min_idle_connects" mapstructure:"min_idle_connects"` // 最小空闲连接.
	IdleTimeout         string        `json:"idle_timeout" yaml:"idle_timeout" mapstructure:"idle_timeout"`                // 空闲链接超时时间
	IdleTimeoutDuration time.Duration `json:"-" yaml:"-" mapstructure:"-"`                                                 // 空闲链接超时时间

	MaxRetries              int           `json:"max_retries" yaml:"max_retries" mapstructure:"max_retries"`                   // 最大重试次数
	MinRetryBackoff         string        `json:"min_retry_backoff" yaml:"min_retry_backoff" mapstructure:"min_retry_backoff"` // 重试策略. 最短重连时间
	MinRetryBackoffDuration time.Duration `json:"-" yaml:"-" mapstructure:"-"`                                                 // 重试策略. 最短重连时间
	MaxRetryBackoff         string        `json:"max_retry_backoff" yaml:"max_retry_backoff" mapstructure:"max_retry_backoff"` // 重试策略. 最大重连时间
	MaxRetryBackoffDuration time.Duration `json:"-" yaml:"-" mapstructure:"-"`                                                 // 重试策略. 最大重连时间

	DialTimeout          string        `json:"dial_timeout" yaml:"dial_timeout" mapstructure:"dial_timeout"`       // 连接超时时间
	DialTimeoutDuration  time.Duration `json:"-" yaml:"-" mapstructure:"-"`                                        // 连接超时时间
	ReadTimeout          string        `json:"read_timeout" yaml:"read_timeout" mapstructure:"dial_timeout"`       // 读超时
	ReadTimeoutDuration  time.Duration `json:"-" yaml:"-" mapstructure:"-"`                                        // 读超时
	WriteTimeout         string        `json:"write_timeout" yaml:"write_timeout" mapstructure:"read_timeout"`     // 写超时
	WriteTimeoutDuration time.Duration `json:"-" yaml:"-" mapstructure:"-"`                                        // 写超时
	ClusterClient        int           `json:"cluster_client" yaml:"cluster_client" mapstructure:"cluster_client"` // 强制集群客户端
}

func (a *Redis) Pretreatment() error {
	var err error
	a.MinRetryBackoffDuration, err = time.ParseDuration(a.MinRetryBackoff)
	if err != nil {
		return err
	}

	a.MaxRetryBackoffDuration, err = time.ParseDuration(a.MaxRetryBackoff)
	if err != nil {
		return err
	}

	a.DialTimeoutDuration, err = time.ParseDuration(a.DialTimeout)
	if err != nil {
		return err
	}

	a.ReadTimeoutDuration, err = time.ParseDuration(a.ReadTimeout)
	if err != nil {
		return err
	}

	a.WriteTimeoutDuration, err = time.ParseDuration(a.WriteTimeout)
	if err != nil {
		return err
	}

	a.IdleTimeoutDuration, err = time.ParseDuration(a.IdleTimeout)
	if err != nil {
		return err
	}
	return nil
}
