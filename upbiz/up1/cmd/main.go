package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/naming"
	"github.com/go-kratos/kratos/pkg/naming/etcd"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"

	"up1/internal/di"
)

func main() {
	flag.Parse()
	log.Init(nil) // debug flag: log.dir={path}
	defer log.Close()
	log.Info("up1 start")
	paladin.Init()
	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}

	etcdV3Conf := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2380"},
		DialTimeout: time.Second * time.Duration(30),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	etcdBuilder, err := etcd.New(&etcdV3Conf)
	if err != nil {
		panic(err)
	}
	ip := "0.0.0.0"
	port := "9002"
	hn, _ := os.Hostname()
	ins := &naming.Instance{
		Zone:     "test",
		Env:      "dev",
		AppID:    "demo.service.up",
		Hostname: hn,
		Addrs: []string{
			"grpc://" + ip + ":" + port,
		},
	}

	cancel, err := etcdBuilder.Register(context.Background(), ins)
	if err != nil {
		panic(err)
	}
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			log.Info("up1 exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
