package dao

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

func NewHttpClient() *blademaster.Client {
	var cfg blademaster.ClientConfig
	var ct paladin.TOML
	if err := paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return nil
	}

	if err := ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return nil
	}
	return blademaster.NewClient(&cfg)
}
