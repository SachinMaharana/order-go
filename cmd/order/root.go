package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	address     = ":8090"
	endpoint    = "127.0.0.1:9091"
	logLevel    = "info"
	redis       = "redis://127.0.0.1:6380/0"
	item        = "http://127.0.0.1:8080"
	rootCommand = &cobra.Command{
		Use:   "order",
		Short: "Simple HTTP order service",
		Run:   runServer,
	}
)

// Execute runs the cobra rootCommand.
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	f := rootCommand.Flags()
	f.StringVarP(&address, "address", "a", address, "Listening Address")
	f.StringVarP(&endpoint, "endpoint", "e", endpoint, "Endpoint for other services to reach order services")
	f.StringVarP(&logLevel, "log-level", "l", logLevel, "log level (debug, info, warn, error), empty or invalid values will fallback to default")
	f.StringVarP(&redis, "redis-address", "r", redis, "redis address to connect to")
	f.StringVarP(&item, "item-address", "i", item, "item service address to query")

}
