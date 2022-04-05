package main

import (
	"fmt"
	"github.com/mabaro3009/example-architecture-go/cmd/command"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	baseCMD := &cobra.Command{
		Use:   "example",
		Short: "base cmd for example app",
	}

	baseCMD.AddCommand(command.Service)

	if err := baseCMD.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
