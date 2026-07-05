package module

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	"starledger/ent"
	reportHandler "starledger/internal/handler/report"
	"starledger/internal/service"
)

type ReportModule struct {
	svc     *service.ReportService
	handler *reportHandler.Handler
}

func NewReportModule() *ReportModule {
	return &ReportModule{}
}

func (m *ReportModule) Name() string {
	return "report"
}

func (m *ReportModule) DisplayName() string {
	return "数据报表"
}

func (m *ReportModule) Description() string {
	return "多维度数据统计分析，包括收支趋势、账单分布、服务器成本分析、任务统计等"
}

func (m *ReportModule) Icon() string {
	return "DataAnalysis"
}

func (m *ReportModule) IsCore() bool {
	return false
}

func (m *ReportModule) Init(client *ent.Client) error {
	m.svc = service.NewReportService(client)
	m.handler = reportHandler.NewHandler(m.svc)
	return nil
}

func (m *ReportModule) RegisterRoutes(group *gin.RouterGroup) {
	reports := group.Group("/reports")
	{
		reports.GET("/overview", m.handler.Overview)
		reports.GET("/monthly-trend", m.handler.MonthlyTrend)
		reports.GET("/bill-type", m.handler.BillTypeDistribution)
		reports.GET("/bill-status", m.handler.BillStatusSummary)
		reports.GET("/server-cost", m.handler.ServerCost)
		reports.GET("/task-stats", m.handler.TaskStats)
	}
}

func (m *ReportModule) RegisterCronJobs(cr *cron.Cron) error {
	return nil
}
