package cmd

import (
	"dot/core"
	"github.com/google/uuid"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"dot/app/v1/deps"
)

var cron = &cobra.Command{
	Use:   "cron",
	Short: "start cron check account activation",
	Run: func(cmd *cobra.Command, args []string) {
		deps := deps.BuildDependency()
		ic := core.NewInternalContext(uuid.NewString())

		go func() {
			log.Info(ic.ToContext(), "cron running...")

			cerr := deps.GetBase().Schlr.Start(ic)
			if cerr != nil {
				log.Error(ic.ToContext(), "cron failed to start", cerr)
			}

			log.Info(ic.ToContext(), "cron finished...")
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		log.Info(ic.ToContext(), "cron shutting down...")
	},
}
