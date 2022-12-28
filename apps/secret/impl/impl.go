package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/lxygwqf9527/demo-cmdb/apps/host"
	"github.com/lxygwqf9527/demo-cmdb/apps/secret"
	"github.com/lxygwqf9527/demo-cmdb/conf"
	"google.golang.org/grpc"
)

var (
	// Service 服务实例
	svr = &impl{}
)

type impl struct {
	db  *sql.DB
	log logger.Logger

	host host.ServiceServer
	secret.UnimplementedServiceServer
}

func (s *impl) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	s.host = app.GetGrpcApp(host.AppName).(host.ServiceServer)
	return nil
}

func (s *impl) Name() string {
	return secret.AppName
}

func (s *impl) Registry(server *grpc.Server) {
	secret.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
