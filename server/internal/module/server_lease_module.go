package module

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"starledger/ent"
	serverhandler "starledger/internal/handler/server"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

type ServerLeaseModule struct {
	svc     *service.ServerLeaseService
	handler *serverhandler.Handler
}

func NewServerLeaseModule() *ServerLeaseModule {
	return &ServerLeaseModule{}
}

func (m *ServerLeaseModule) Name() string {
	return "server_lease"
}

func (m *ServerLeaseModule) DisplayName() string {
	return "服务器租赁"
}

func (m *ServerLeaseModule) Description() string {
	return "管理服务器租赁信息，包括服务商、配置、费用、到期提醒、续租等"
}

func (m *ServerLeaseModule) Icon() string {
	return "Monitor"
}

func (m *ServerLeaseModule) IsCore() bool {
	return false
}

func (m *ServerLeaseModule) Init(client *ent.Client) error {
	m.svc = service.NewServerLeaseService(client)
	m.handler = serverhandler.NewHandler(m.svc)
	return nil
}

func (m *ServerLeaseModule) RegisterRoutes(group *gin.RouterGroup) {
	servers := group.Group("/servers")
	{
		servers.GET("", m.handler.List)
		servers.GET("/expiring", m.handler.Expiring)
		servers.GET("/:id", m.handler.Get)
		servers.POST("", m.handler.Create)
		servers.PUT("/:id", m.handler.Update)
		servers.DELETE("/:id", m.handler.Delete)
		servers.PUT("/:id/renew", m.handler.Renew)
	}
}

func (m *ServerLeaseModule) RegisterCronJobs(cr *cron.Cron) error {
	_, err := cr.AddFunc("0 0 0 * * *", func() {
		if err := m.svc.CheckExpiring(context.Background()); err != nil {
			pkg.Logger.Error("server_lease: check expiring failed", zap.Error(err))
		}
	})
	return err
}
