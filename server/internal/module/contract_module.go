package module

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"starledger/ent"
	contracthandler "starledger/internal/handler/contract"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

type ContractModule struct {
	svc     *service.ContractService
	handler *contracthandler.Handler
}

func NewContractModule() *ContractModule {
	return &ContractModule{}
}

func (m *ContractModule) Name() string {
	return "contract"
}

func (m *ContractModule) DisplayName() string {
	return "合同管理"
}

func (m *ContractModule) Description() string {
	return "管理企业合同信息，包括合同创建、签署、到期跟踪等"
}

func (m *ContractModule) Icon() string {
	return "Tickets"
}

func (m *ContractModule) IsCore() bool {
	return false
}

func (m *ContractModule) Init(client *ent.Client) error {
	m.svc = service.NewContractService(client)
	m.handler = contracthandler.NewHandler(m.svc)
	return nil
}

func (m *ContractModule) RegisterRoutes(group *gin.RouterGroup) {
	contracts := group.Group("/contracts")
	{
		contracts.GET("", m.handler.List)
		contracts.GET("/:id", m.handler.Get)
		contracts.POST("", m.handler.Create)
		contracts.PUT("/:id", m.handler.Update)
		contracts.DELETE("/:id", m.handler.Delete)
	}
}

func (m *ContractModule) RegisterCronJobs(cr *cron.Cron) error {
	_, err := cr.AddFunc("0 0 0 * * *", func() {
		if err := m.svc.CheckExpiring(context.Background()); err != nil {
			pkg.Logger.Error("contract: check expiring failed", zap.Error(err))
		}
	})
	return err
}
