package command

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/mabaro3009/example-architecture-go/service"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var Service = &cobra.Command{
	Use:   "service",
	Short: "Service server",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := service.Config{}
		if err := envconfig.Process("example", &conf); err != nil {
			return err
		}

		srv, err := service.NewService(conf)
		if err != nil {
			return err
		}

		fmt.Printf("listening on http://%s\n", conf.ListenAddress)

		go srv.ListenAndServe()

		// Waiting for an OS signal cancellation
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		srv.Shutdown()
		return nil
	},
}
