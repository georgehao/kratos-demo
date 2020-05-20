// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"up1/internal/dao"
	"up1/internal/service"
	"up1/internal/server/grpc"
	"up1/internal/server/http"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, grpc.New, NewApp))
}
