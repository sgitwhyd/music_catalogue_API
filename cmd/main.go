package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/sgitwhyd/music-catalogue/internal/configs"
	"github.com/sgitwhyd/music-catalogue/pkg/internalsql"
)

func main() {
	config, err := configs.Init("./", "env", ".env")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Info("failed load config")
	}
	log.Infoln("load config success")

	_, err = internalsql.Connect(config.DatabaseURL)
	if err != nil {
	log.WithFields(log.Fields{
			"error": err.Error(),
		}).Info("failed connecting database")
	}
	log.Info("database connected")

	if config.ENV == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	route := gin.Default()
	route.Run(config.PORT)
}