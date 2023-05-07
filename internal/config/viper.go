package config

import (
	"github.com/spf13/viper"
)

var Global Wrap

const (
	Local   = 1 // 本地环境
	China   = 2 // 国内环境
	Foreign = 3 // 国外环境
)

func NewViper(path string) (*viper.Viper, Wrap, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		return nil, Wrap{}, err
	}

	var c Wrap
	if err = v.Unmarshal(&c); err != nil {
		return nil, Wrap{}, err
	}

	err = c.Pretreatment()
	if err != nil {
		return nil, Wrap{}, err
	}
	return v, c, nil
}

func Init() error {
	var err error
	_, Global, err = NewViper(FilePath())
	if err != nil {
		return err
	}
	return nil
}

// IsChina 是国内环境
func IsChina() bool {
	return Global.App.Development == China
}

// IsLocal 是本地环境
func IsLocal() bool {
	return Global.App.Development == Local
}

// IsForeign 是国外环境
func IsForeign() bool {
	return Global.App.Development == Foreign
}
