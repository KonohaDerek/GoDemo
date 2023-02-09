package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"k-derek.dev/demo/cqrs/cmd"
)

var rootCmd = &cobra.Command{Use: "server"}

// @title Gin swagger
// @version 1.0
// @description Gin swagger

// @contact.name miyo
// @contact.url https://DerekChenDev@bitbucket.org/xspinach/man_backend.git

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8087
// schemes http
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	rootCmd.AddCommand(cmd.ServerCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
