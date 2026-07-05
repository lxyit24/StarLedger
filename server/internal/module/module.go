package module

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"starledger/ent"
	"starledger/internal/middleware"
	"starledger/internal/pkg"
)

type Module interface {
	Name() string
	DisplayName() string
	Description() string
	Icon() string
	IsCore() bool
	Init(client *ent.Client) error
	RegisterRoutes(group *gin.RouterGroup)
	RegisterCronJobs(cr *cron.Cron) error
}

// ModuleInfo holds module metadata for marketplace display.
type ModuleInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IsCore      bool   `json:"is_core"`
	Enabled     bool   `json:"enabled"`
}

type Registry struct {
	mods   []Module
	client *ent.Client
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(m Module) {
	r.mods = append(r.mods, m)
	pkg.Logger.Info("module registered", zap.String("module", m.Name()))
}

func (r *Registry) InitAll(client *ent.Client) error {
	r.client = client
	for _, m := range r.mods {
		pkg.Logger.Info("initializing module", zap.String("module", m.Name()))
		if err := m.Init(client); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) RegisterAllRoutes(group *gin.RouterGroup) {
	for _, m := range r.mods {
		pkg.Logger.Info("registering routes for module", zap.String("module", m.Name()))
		// Wrap module routes with access control middleware
		modGroup := group.Group("")
		if r.client != nil && !m.IsCore() {
			modGroup.Use(middleware.ModuleAccess(r.client, m.Name()))
		}
		m.RegisterRoutes(modGroup)
	}
}

func (r *Registry) RegisterAllCronJobs(cr *cron.Cron) error {
	for _, m := range r.mods {
		pkg.Logger.Info("registering cron jobs for module", zap.String("module", m.Name()))
		if err := m.RegisterCronJobs(cr); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) ListModules() []string {
	names := make([]string, len(r.mods))
	for i, m := range r.mods {
		names[i] = m.Name()
	}
	return names
}

// ListModuleInfos returns metadata for all registered modules.
func (r *Registry) ListModuleInfos() []ModuleInfo {
	infos := make([]ModuleInfo, len(r.mods))
	for i, m := range r.mods {
		infos[i] = ModuleInfo{
			Name:        m.Name(),
			DisplayName: m.DisplayName(),
			Description: m.Description(),
			Icon:        m.Icon(),
			IsCore:      m.IsCore(),
		}
	}
	return infos
}

// GetModule returns a module by name.
func (r *Registry) GetModule(name string) (Module, bool) {
	for _, m := range r.mods {
		if m.Name() == name {
			return m, true
		}
	}
	return nil, false
}
