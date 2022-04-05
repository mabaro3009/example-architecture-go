package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/mabaro3009/example-architecture-go/service"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	conf := service.Config{}
	if err := envconfig.Process("example", &conf); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	srv, err := service.NewService(conf)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("listening on http://localhost%s\n", conf.ListenAddress)

	go srv.ListenAndServe()

	// Waiting for an OS signal cancellation
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	srv.Shutdown()
	return
}
