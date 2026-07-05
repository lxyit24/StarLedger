package module

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	"starledger/ent"
	auditHandler "starledger/internal/handler/audit"
	"starledger/internal/service"
)

type AuditModule struct {
	svc     *service.AuditLogService
	handler *auditHandler.Handler
}

func NewAuditModule() *AuditModule {
	return &AuditModule{}
}

func (m *AuditModule) Name() string {
	return "audit"
}

func (m *AuditModule) DisplayName() string {
	return "审计日志"
}

func (m *AuditModule) Description() string {
	return "记录所有关键操作日志，支持查询和追溯，保障数据安全合规"
}

func (m *AuditModule) Icon() string {
	return "Notebook"
}

func (m *AuditModule) IsCore() bool {
	return false
}

func (m *AuditModule) Init(client *ent.Client) error {
	m.svc = service.NewAuditLogService(client)
	m.handler = auditHandler.NewHandler(m.svc)
	return nil
}

func (m *AuditModule) RegisterRoutes(group *gin.RouterGroup) {
	audit := group.Group("/audit")
	{
		audit.GET("/logs", m.handler.List)
		audit.POST("/logs", m.handler.Create)
	}
}

func (m *AuditModule) RegisterCronJobs(cr *cron.Cron) error {
	return nil
}
