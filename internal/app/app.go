package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Woland-prj/dilemator/config"
	"github.com/Woland-prj/dilemator/internal/domain/errors/berrors"
	"github.com/Woland-prj/dilemator/internal/router/handlers/api/dilemma_router"
	"github.com/Woland-prj/dilemator/internal/router/handlers/api/sessions_router"
	"github.com/Woland-prj/dilemator/internal/router/handlers/api/users_router"
	"github.com/Woland-prj/dilemator/internal/router/handlers/pages_router"
	"github.com/Woland-prj/dilemator/internal/router/managers"
	router_setup "github.com/Woland-prj/dilemator/internal/router/setup"
	"github.com/Woland-prj/dilemator/internal/services/factory"
	"github.com/Woland-prj/dilemator/internal/services/sessions_service"
	"github.com/Woland-prj/dilemator/pkg/httpserver"
	"github.com/Woland-prj/dilemator/pkg/logger"
)

func Run(cfg *config.Config) {
	log, err := logger.New(cfg.Log.Env, cfg.Log.Multiple, cfg.Log.File)
	if err != nil {
		panic(fmt.Errorf("app - Run - logger.New: %w", err))
	}

	// Creating a core services factory
	f, err := factory.NewServiceFactory(cfg, log)
	if err != nil {
		log.Error(fmt.Errorf("app - Run - factory.NewServiceFactory: %w", err).Error())
		os.Exit(1)
	}

	// SessionService for cookieManager
	ss, err := factory.InstantiateService[sessions_service.SessionService](f)
	if err != nil {
		log.Error(fmt.Errorf("app - Run - factory.InstantiateService: %w", err).Error())
		os.Exit(1)
	}

	// Create server
	httpServer := httpserver.New(
		httpserver.Port(cfg.HTTP.Port),
		httpserver.Prefork(cfg.HTTP.UsePreforkMode),
		httpserver.BodyLimitMb(cfg.HTTP.BodyLimitMb),
	)

	// Create cookie manager
	cookieManager := managers.NewCookieManager(ss, log, getCookieConf(cfg))

	// Register routes
	mustRegisterRoutes(f, httpServer, cookieManager, cfg, log)

	// Start server
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal:", slog.Any("signal", s))
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}

func getCookieConf(cfg *config.Config) *managers.CookieConfig {
	return &managers.CookieConfig{
		CookieName: cfg.HTTP.CookieName,
		MaxAgeDays: cfg.HTTP.CookieMaxAgeDays,
		Secure:     cfg.HTTP.CookieSecure,
		SameSite:   cfg.HTTP.CookieSameSite,
	}
}

func mustRegisterRoutes(
	f *factory.ServiceFactory,
	httpServer *httpserver.Server,
	cm *managers.CookieManager,
	cfg *config.Config,
	log logger.Interface,
) {
	const op = "app - mustRegisterRoutes"

	apiRouter, componentsRouter := router_setup.NewRouter(httpServer.App, cfg, log)

	err := users_router.Register(apiRouter, componentsRouter, f, cm, log)
	if err != nil {
		log.Error(berrors.FromErr(op, err).Error())
		os.Exit(1)
	}

	err = dilemma_router.Register(apiRouter, componentsRouter, f, cm, log)
	if err != nil {
		log.Error(berrors.FromErr(op, err).Error())
		os.Exit(1)
	}

	err = sessions_router.Register(apiRouter, f, cm, log)
	if err != nil {
		log.Error(berrors.FromErr(op, err).Error())
		os.Exit(1)
	}

	err = pages_router.Register(httpServer.App, log, cfg.App.TgBotName)
	if err != nil {
		log.Error(berrors.FromErr(op, err).Error())
		os.Exit(1)
	}
}
