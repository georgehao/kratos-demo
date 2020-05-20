package service

import (
	"context"
	"github.com/go-kratos/kratos/pkg/log"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
	pb "zbbiz/api"
	"zbbiz/internal/dao"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.UpServer), new(*Service)))

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

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}

func (s *Service) CheckRole(ctx context.Context, req *pb.CheckUpReq) (resp *pb.CheckUpResp, err error) {
	log.Infoc(ctx, "check role:%d", req.Role)
	resp, err = s.dao.CheckRole(ctx, req)
	return
}
