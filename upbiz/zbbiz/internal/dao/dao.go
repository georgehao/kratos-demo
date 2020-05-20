package dao

import (
	"context"
	"github.com/go-kratos/kratos/pkg/cache/memcache"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/sync/pipeline/fanout"
	xtime "github.com/go-kratos/kratos/pkg/time"
	"github.com/google/wire"
	"time"
	api "zbbiz/api"
)

var Provider = wire.NewSet(New, NewDB, NewRedis, NewMC)

//go:generate kratos tool genbts
// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	CheckRole(ctx context.Context, req *api.CheckUpReq) (resp *api.CheckUpResp, err error)
}

// dao dao.
type dao struct {
	db         *sql.DB
	redis      *redis.Redis
	mc         *memcache.Memcache
	cache      *fanout.Fanout
	demoExpire int32
	rpcClient  api.UpClient
	httpClient blademaster.Client
}

// New new a dao and return.
func New(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d Dao, cf func(), err error) {
	return newDao(r, mc, db)
}

func newDao(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d *dao, cf func(), err error) {
	var cfg struct {
		DemoExpire xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}

	var grpcClient api.UpClient
	grpcClient, err = NewGrpcClient()
	if err != nil {
		return
	}

	d = &dao{
		db:         db,
		redis:      r,
		mc:         mc,
		cache:      fanout.New("cache"),
		demoExpire: int32(time.Duration(cfg.DemoExpire) / time.Second),
		rpcClient:  grpcClient,
		httpClient: *NewHttpClient(),
	}
	cf = d.Close
	return
}

// Close close the resource.
func (d *dao) Close() {
	d.cache.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return nil
}

type roleResult struct {
	Yes bool `json:"yes"`
}

func (d *dao) CheckRole(ctx context.Context, req *api.CheckUpReq) (resp *api.CheckUpResp, err error) {
	resp, err = d.rpcClient.CheckRole(ctx, req)
	return
}

//func (d *dao) CheckRole(ctx context.Context, req *api.CheckUpReq) (resp *api.CheckUpResp, err error) {
//	params := url.Values{}
//	params.Set("roleId", fmt.Sprintf("%d", req.Role))
//
//	var resp1 struct {
//		Code int        `json:"code"`
//		Data roleResult `json:"data"`
//	}
//
//	upUrl := "http://127.0.0.1:8000/up/check_role"
//	if err = d.httpClient.Get(context.Background(), upUrl, "", params, &resp1); err != nil {
//		return nil, err
//	}
//
//	if resp1.Code != 0 {
//		err = errors.Errorf("up url(%s) res(%+v) err(%+v)", upUrl+"?"+params.Encode(), resp, ecode.Int(resp1.Code))
//		return
//	}
//
//	resp = &api.CheckUpResp{
//		Yes: resp1.Data.Yes,
//	}
//	return
//}
