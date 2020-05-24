package service

import (
	"context"
	pb "fanout/api"
	"fanout/internal/dao"
	"github.com/go-kratos/kratos/pkg/conf/paladin"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.DemoServer), new(*Service)))

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

// New new a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

func (s *Service) UpdateRole(ctx context.Context, req *pb.UpdateRoleReq) (resp *pb.UpdateRoleResp, err error) {
	resp = &pb.UpdateRoleResp{}
	_, err = s.dao.UpdateRole(ctx, req.Role)
	if err != nil {
		return resp, err
	}
	resp.Yes = true

	// 刷新redis
	err = s.dao.FanoutDo(ctx, func(c context.Context) {
		s.dao.UpdateRoleRedis(c, 100, req.Role)
	})
	return
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}
