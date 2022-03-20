package bootstrap

import (
	"context"

	"github.com/Cray-HPE/cray-nls/api/controllers"
	"github.com/Cray-HPE/cray-nls/api/middlewares"
	"github.com/Cray-HPE/cray-nls/api/routes"
	"github.com/Cray-HPE/cray-nls/services"
	"github.com/Cray-HPE/cray-nls/utils"
	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	controllers.Module,
	routes.Module,
	utils.Module,
	services.Module,
	middlewares.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler utils.RequestHandler,
	routes routes.Routes,
	env utils.Env,
	logger utils.Logger,
	middlewares middlewares.Middlewares,
) {

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting Application")
			logger.Info("---------------------")
			logger.Info("------- CLEAN -------")
			logger.Info("---------------------")

			go func() {
				middlewares.Setup()
				routes.Setup()
				host := "0.0.0.0"
				if env.Environment == "development" {
					host = "127.0.0.1"
				}
				handler.Gin.Run(host + ":" + env.ServerPort)
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("Stopping Application")
			return nil
		},
	})
}
