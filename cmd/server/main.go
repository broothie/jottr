package main

import (
	"fmt"
	"net/http"

	"github.com/broothie/jottr/config"
	"github.com/broothie/jottr/logger"
	"github.com/broothie/jottr/server"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("jottr", &cfg); err != nil {
		panic(err)
	}

	var loggerConfig []logger.Configurer
	if cfg.Environment != "production" {
		loggerConfig = append(loggerConfig, logger.UseHumanFormat())
	}

	log := logger.New(loggerConfig...)
	defer log.Close()

	handler, err := server.New(cfg, log)
	if err != nil {
		log.Err(err, "failed to create server")
		return
	}

	log.Info("server running", logger.Field("config", fmt.Sprintf("%+v", cfg)))
	log.Err(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), handler), "server stopped")
}
