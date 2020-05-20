package http

import (
	"errors"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/ecode"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	pb "zbbiz/api"
	"zbbiz/internal/model"
)

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
	g := e.Group("/zbbiz")
	{
		g.GET("/check_role", checkRole)
	}
}

func ping(ctx *bm.Context) {
}

type roleParam struct {
	RoleId int32 `form:"roleId"`
}

func checkRole(context *bm.Context) {
	var input roleParam
	if err := context.Bind(&input); err != nil {
		context.JSON(nil, ecode.RequestErr)
		return
	}

	resp, err := svc.CheckRole(context, &pb.CheckUpReq{Role: input.RoleId})
	if err != nil {
		context.JSON(nil, err)
		return
	}

	if resp == nil {
		context.JSON(nil, errors.New("resp is nil"))
		return
	}

	cr := model.CheckRole{
		Yes: resp.Yes,
	}
	context.JSON(cr, nil)
}
