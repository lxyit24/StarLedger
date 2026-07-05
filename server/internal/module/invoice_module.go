package module

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	"starledger/ent"
	invoiceHandler "starledger/internal/handler/invoice"
	"starledger/internal/service"
)

type InvoiceModule struct {
	svc     *service.InvoiceService
	handler *invoiceHandler.Handler
}

func NewInvoiceModule() *InvoiceModule {
	return &InvoiceModule{}
}

func (m *InvoiceModule) Name() string {
	return "invoice"
}

func (m *InvoiceModule) DisplayName() string {
	return "发票管理"
}

func (m *InvoiceModule) Description() string {
	return "发票开具、管理、查询，支持增值税普通/专用发票，关联账单管理"
}

func (m *InvoiceModule) Icon() string {
	return "Ticket"
}

func (m *InvoiceModule) IsCore() bool {
	return false
}

func (m *InvoiceModule) Init(client *ent.Client) error {
	m.svc = service.NewInvoiceService(client)
	m.handler = invoiceHandler.NewHandler(m.svc)
	return nil
}

func (m *InvoiceModule) RegisterRoutes(group *gin.RouterGroup) {
	invoices := group.Group("/invoices")
	{
		invoices.GET("", m.handler.List)
		invoices.GET("/summary", m.handler.Summary)
		invoices.GET("/:id", m.handler.Get)
		invoices.POST("", m.handler.Create)
		invoices.PUT("/:id", m.handler.Update)
		invoices.PUT("/:id/issue", m.handler.Issue)
		invoices.PUT("/:id/cancel", m.handler.Cancel)
		invoices.DELETE("/:id", m.handler.Delete)
	}
}

func (m *InvoiceModule) RegisterCronJobs(cr *cron.Cron) error {
	return nil
}
