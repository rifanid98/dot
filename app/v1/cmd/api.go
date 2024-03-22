package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.elastic.co/apm/module/apmechov4"

	"dot/app/v1/deps"
	"dot/interface/v1/general/common"

	routerV0 "dot/interface/v0/router"
	appMiddleware "dot/interface/v1/general/middleware"
	routerV1 "dot/interface/v1/general/router"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start api service",
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()
		e.HideBanner = false
		e.HidePort = false
		e.Validator = common.NewValidator()

		e.Use(apmechov4.Middleware())
		e.Use(appMiddleware.ServiceTrackerID)
		e.Use(appMiddleware.ServiceRequestTime)
		e.Use(echoMiddleware.RemoveTrailingSlash())
		e.Use(appMiddleware.Recover())
		e.Use(appMiddleware.CORS())

		deps := deps.BuildDependency()
		routerV0.Register(e, deps)
		routerV1.Register(e, deps)

		go func() {
			if err := e.Start(fmt.Sprintf(":%v", deps.GetBase().Cfg.Port)); err != nil {
				panic(err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			panic(err)
		}
	},
}

func ExecuteApiCommand() {
	apiCmd.AddCommand(cron)
	if err := apiCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
