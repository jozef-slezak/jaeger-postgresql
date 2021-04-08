package main

import (
        "os"
 
	"github.com/hashicorp/go-hclog"
	"jaeger-postgresql/pgstore"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc/shared"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "jaeger-postgresql",
		Level: hclog.Warn, // Jaeger only captures >= Warn, so don't bother logging below Warn
	})

	var store shared.StoragePlugin
	var closeStore func() error
	var err error

	conf := pgstore.Configuration{}

	store, closeStore, err = pgstore.NewStore(&conf, logger)

	grpc.Serve(store)

	if err = closeStore(); err != nil {
		logger.Error("failed to close store", "error", err)
		os.Exit(1)
	}
}
