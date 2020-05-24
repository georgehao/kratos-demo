// Code generated by protoc-gen-bm v0.1, DO NOT EDIT.
// source: api.proto

/*
Package api is a generated blademaster stub package.
This code was generated with kratos/tool/protobuf/protoc-gen-bm v0.1.

package 命名使用 {appid}.{version} 的方式, version 形如 v1, v2 ..

It is generated from these files:
	api.proto
*/
package api

import (
	"context"

	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
)
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

// to suppressed 'imported but not used warning'
var _ *bm.Context
var _ context.Context
var _ binding.StructValidator

var PathDemoPing = "/fanout.service.v1.Demo/Ping"
var PathDemoUpdateRole = "/fanout.service.v1.Demo/UpdateRole"

// DemoBMServer is the server API for Demo service.
type DemoBMServer interface {
	Ping(ctx context.Context, req *google_protobuf.Empty) (resp *google_protobuf.Empty, err error)

	UpdateRole(ctx context.Context, req *UpdateRoleReq) (resp *UpdateRoleResp, err error)
}

var DemoSvc DemoBMServer

func demoPing(c *bm.Context) {
	p := new(google_protobuf.Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.Ping(c, p)
	c.JSON(resp, err)
}

func demoUpdateRole(c *bm.Context) {
	p := new(UpdateRoleReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := DemoSvc.UpdateRole(c, p)
	c.JSON(resp, err)
}

// RegisterDemoBMServer Register the blademaster route
func RegisterDemoBMServer(e *bm.Engine, server DemoBMServer) {
	DemoSvc = server
	e.GET("/fanout.service.v1.Demo/Ping", demoPing)
	e.GET("/fanout.service.v1.Demo/UpdateRole", demoUpdateRole)
}
