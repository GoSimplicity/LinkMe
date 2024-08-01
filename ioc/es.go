package ioc

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
)

// InitES 初始化elasticsearch
func InitES() *elasticsearch.Client {
	addr := viper.GetString("es.addr")
	cfg := elasticsearch.Config{
		Addresses: []string{
			addr,
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return client
}
