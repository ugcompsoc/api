// Credit
// Thanks to TUD NetSoc for this brilliant entry point.
// https://github.com/netsoc/webspaced

package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/nuigcompsoc/api/internal/config"
	"github.com/nuigcompsoc/api/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var version = "dev"
var srv *server.Server

func init() {
	// Config defaults
	viper.SetDefault("log_level", log.InfoLevel)

	viper.SetDefault("timeouts.startup", 30*time.Second)
	viper.SetDefault("timeouts.shutdown", 30*time.Second)

	viper.SetDefault("http.listen_address", ":80")
	viper.SetDefault("http.cors.allowed_origins", []string{"*"})

	// Config file loading
	viper.SetConfigType("yaml")
	viper.SetConfigName("api")
	viper.AddConfigPath("/run/config")
	viper.AddConfigPath(".")

	// Config from environment
	viper.SetEnvPrefix("ENV")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Config from flags
	pflag.StringP("log_level", "l", "info", "log level")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.WithError(err).Fatal("Failed to bind pflags to config")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.WithField("err", err).Warn("Failed to read configuration")
	}
}

func reload() {
	if srv != nil {
		stop()
		srv = nil
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg, config.DecoderOptions); err != nil {
		log.WithField("err", err).Fatal("Failed to parse configuration")
	}

	log.SetLevel(cfg.LogLevel)
	cJSON, err := json.Marshal(cfg)
	if err != nil {
		log.WithError(err).Fatal("Failed to encode config as JSON")
	}
	log.WithField("config", string(cJSON)).Debug("Got config")

	srv = server.NewServer(cfg)

	log.Info("Starting server")
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Startup)
		defer cancel()
		if err := srv.Start(ctx); err != nil {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()
}

func stop() {
	log.Info("Stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), srv.Config.Timeouts.Shutdown)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.WithError(err).Fatal("Failed to stop server")
	}

	log.Info("Stopped server sucessfully")
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.WithField("file", e.Name).Info("Config changed, reloading")
		reload()
	})
	viper.WatchConfig()
	reload()

	<-sigs
	stop()
}
