package main

import (
	"flag"
	"github.com/hashicorp/go-hclog"
	"github.com/innovatrics/jaeger-postgresql/pgstore"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Warn,
		Name:       "jaeger-postgresql",
		JSONFormat: true,
	})

	var configPath string
	flag.StringVar(&configPath, "config", "", "A path to the plugin's configuration file")
	flag.Parse()

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if configPath != "" {
		v.SetConfigFile(configPath)
	}

	if configPath != "" {
		err := v.ReadInConfig()
		if err != nil {
			logger.Error("failed to parse configuration file", "error", err)
			os.Exit(1)
		}
	}

	var conf pgstore.Configuration
	conf.InitFromViper(v)

	store, f, err := pgstore.NewStore(&conf, logger)
	if err != nil {
		logger.Error("failed to create postgresql store", "error", err)
		err := f()
		if err != nil {
			logger.Error("failed to close postgresql store", "error", err)
		}
		os.Exit(1)
	}

	grpc.Serve(store)
}
