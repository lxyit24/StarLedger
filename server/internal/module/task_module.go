package module

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	"starledger/ent"
	taskhandler "starledger/internal/handler/task"
	"starledger/internal/service"
)

type TaskModule struct {
	svc     *service.TaskService
	handler *taskhandler.Handler
}

func NewTaskModule() *TaskModule {
	return &TaskModule{}
}

func (m *TaskModule) Name() string {
	return "task"
}

func (m *TaskModule) DisplayName() string {
	return "任务协作"
}

func (m *TaskModule) Description() string {
	return "团队协作任务管理，包括任务创建、分配、进度跟踪等"
}

func (m *TaskModule) Icon() string {
	return "List"
}

func (m *TaskModule) IsCore() bool {
	return false
}

func (m *TaskModule) Init(client *ent.Client) error {
	m.svc = service.NewTaskService(client)
	m.handler = taskhandler.NewHandler(m.svc)
	return nil
}

func (m *TaskModule) RegisterRoutes(group *gin.RouterGroup) {
	tasks := group.Group("/tasks")
	{
		tasks.GET("", m.handler.List)
		tasks.GET("/my", m.handler.MyTasks)
		tasks.GET("/:id", m.handler.Get)
		tasks.POST("", m.handler.Create)
		tasks.PUT("/:id", m.handler.Update)
		tasks.PUT("/:id/assign", m.handler.Assign)
		tasks.DELETE("/:id", m.handler.Delete)
	}
}

func (m *TaskModule) RegisterCronJobs(cr *cron.Cron) error {
	return nil
}
