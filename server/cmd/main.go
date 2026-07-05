package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"starledger/ent"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"starledger/internal/config"
	"starledger/internal/handler/auth"
	"starledger/internal/handler/market"
	"starledger/internal/handler/tenant"
	"starledger/internal/handler/user"
	"starledger/internal/middleware"
	"starledger/internal/module"
	"starledger/internal/pkg"
	"starledger/internal/service"
)

func main() {
	cfgPath := "config.yaml"
	if envPath := os.Getenv("STARLEDGER_CONFIG"); envPath != "" {
		cfgPath = envPath
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	pkg.InitLogger()
	defer pkg.Logger.Sync()

	db, err := sql.Open("sqlite", cfg.Database.DSN())
	if err != nil {
		pkg.Logger.Fatal("failed to connect database", zap.Error(err))
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		pkg.Logger.Fatal("failed to enable foreign keys", zap.Error(err))
	}
	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		pkg.Logger.Fatal("failed to create schema", zap.Error(err))
	}
	pkg.Logger.Info("database schema migrated successfully")

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authSvc := service.NewAuthService(client)
	userSvc := service.NewUserService(client)
	tenantSvc := service.NewTenantService(client)
	roleSvc := service.NewRoleService(client)
	marketSvc := service.NewMarketService(client)

	authHandler := auth.NewHandler(authSvc, userSvc)
	tenantHandler := tenant.NewHandler(tenantSvc)
	userHandler := user.NewHandler(userSvc)
	roleHandler := user.NewRoleHandler(roleSvc)

	registry := module.NewRegistry()
	registry.Register(module.NewServerLeaseModule())
	registry.Register(module.NewBillingModule())
	registry.Register(module.NewContractModule())
	registry.Register(module.NewTaskModule())

	if err := registry.InitAll(client); err != nil {
		pkg.Logger.Fatal("failed to init modules", zap.Error(err))
	}
	pkg.Logger.Info("all modules initialized", zap.Strings("modules", registry.ListModules()))

	marketHandler := market.NewHandler(marketSvc, registry)

	v1 := r.Group("/api/v1")
	{
		public := v1.Group("/auth")
		{
			public.POST("/login", authHandler.Login)
			public.POST("/register", authHandler.Register)
		}

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(), middleware.TenantIsolation())
		{
			protected.GET("/auth/profile", authHandler.Profile)
			protected.PUT("/auth/password", authHandler.ChangePassword)

			tenants := protected.Group("/tenants")
			{
				tenants.GET("", tenantHandler.List)
				tenants.POST("", tenantHandler.Create)
				tenants.PUT("/:id", tenantHandler.Update)
				tenants.DELETE("/:id", tenantHandler.Delete)
			}

			users := protected.Group("/users")
			{
				users.GET("", userHandler.List)
				users.POST("", userHandler.Create)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Delete)
			}

			roles := protected.Group("/roles")
			{
				roles.GET("", roleHandler.List)
				roles.POST("", roleHandler.Create)
				roles.PUT("/:id", roleHandler.Update)
			}

			market := protected.Group("/market")
			{
				market.GET("/modules", marketHandler.ListModules)
				market.PUT("/modules/:name/enable", marketHandler.EnableModule)
				market.PUT("/modules/:name/disable", marketHandler.DisableModule)
			}

			registry.RegisterAllRoutes(protected)
		}
	}

	cr := cron.New(cron.WithSeconds())
	if err := registry.RegisterAllCronJobs(cr); err != nil {
		pkg.Logger.Error("failed to register cron jobs", zap.Error(err))
	}
	cr.Start()
	defer cr.Stop()

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		pkg.Logger.Info(fmt.Sprintf("server starting on %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkg.Logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	pkg.Logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		pkg.Logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	pkg.Logger.Info("server exited")
}
