package module

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"starledger/ent"
	billinghandler "starledger/internal/handler/billing"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

type BillingModule struct {
	svc     *service.BillService
	handler *billinghandler.Handler
}

func NewBillingModule() *BillingModule {
	return &BillingModule{}
}

func (m *BillingModule) Name() string {
	return "billing"
}

func (m *BillingModule) DisplayName() string {
	return "账单管理"
}

func (m *BillingModule) Description() string {
	return "统一管理所有业务账单，包括生成、支付、逾期检查、汇总统计等"
}

func (m *BillingModule) Icon() string {
	return "Document"
}

func (m *BillingModule) IsCore() bool {
	return true
}

func (m *BillingModule) Init(client *ent.Client) error {
	m.svc = service.NewBillService(client)
	m.handler = billinghandler.NewHandler(m.svc)
	return nil
}

func (m *BillingModule) RegisterRoutes(group *gin.RouterGroup) {
	bills := group.Group("/bills")
	{
		bills.GET("", m.handler.List)
		bills.GET("/summary", m.handler.Summary)
		bills.GET("/overdue", m.handler.Overdue)
		bills.GET("/:id", m.handler.Get)
		bills.POST("", m.handler.Create)
		bills.PUT("/:id/pay", m.handler.Pay)
		bills.PUT("/:id/cancel", m.handler.Cancel)
	}
}

func (m *BillingModule) RegisterCronJobs(cr *cron.Cron) error {
	_, err := cr.AddFunc("0 0 1 * * *", func() {
		if err := m.svc.CheckOverdue(context.Background()); err != nil {
			pkg.Logger.Error("billing: check overdue failed", zap.Error(err))
		}
	})
	return err
}
