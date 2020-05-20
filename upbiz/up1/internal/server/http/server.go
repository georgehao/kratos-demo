package http

import (
	"errors"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"net/http"
	pb "up1/api"
)

//var svc pb.DemoServer
var svc pb.UpServer

// New new a bm server.
func New(s pb.UpServer) (engine *bm.Engine, err error) {
	var (
		cfg bm.ServerConfig
		ct  paladin.TOML
	)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	svc = s
	engine = bm.DefaultServer(&cfg)
	pb.RegisterUpBMServer(engine, s)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/up")
	{
		g.GET("/check_role", checkRole)
	}
}

type roleParam struct {
	RoleId int32 `form:"roleId"`
}

type roleResult struct {
	Yes bool `json:"yes"`
}

func checkRole(context *bm.Context) {
	var input roleParam
	if err := context.Bind(&input); err != nil {
		context.JSON(nil, errors.New("client param error"))
		return
	}

	resp, err := svc.CheckRole(context, &pb.CheckUpReq{Role: input.RoleId})
	if err != nil {
		context.JSON(nil, err)
		return
	}
	res := roleResult{Yes: resp.Yes}
	context.JSON(res, nil)
}

func ping(ctx *bm.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
