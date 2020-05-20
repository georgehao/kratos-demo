package dao

import (
	"context"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/naming/etcd"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden/resolver"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"time"
	"zbbiz/api"
)

//const target = "direct://default/127.0.0.1:9002"

const AppID = "demo.service.up" // NOTE: example

func init() {
	etcdV3Conf := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2380"},
		DialTimeout: time.Second * time.Duration(30),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}
	resolver.Register(etcd.Builder(&etcdV3Conf))
}

func NewGrpcClient() (grpcClient api.UpClient, err error) {
	cfg := &warden.ClientConfig{}
	var ct paladin.TOML
	if err := paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
		return nil, err
	}

	if err := ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return nil, err
	}

	client := warden.NewClient(cfg)
	//cc, err := client.Dial(context.Background(), target)
	//cc, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	cc, err := client.Dial(context.Background(), "etcd://default/"+AppID)
	if err != nil {
		return nil, err
	}
	return api.NewUpClient(cc), nil
}
