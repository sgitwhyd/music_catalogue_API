package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/sgitwhyd/music-catalogue/internal/configs"
	"github.com/sgitwhyd/music-catalogue/internal/models"
	"github.com/sgitwhyd/music-catalogue/pkg/internalsql"
)

func main() {
	config, err := configs.Init("./", "env", ".env")
	if err != nil {
		log.Error().Err(err).Msgf("load config error %v", err)
	}
	log.Info().Msg("load config success")

	db, err := internalsql.Connect(config.DatabaseURL)
	if err != nil {
		log.Error().Err(err).Msgf("error db connection %v", err)
	}
	log.Info().Msg("database connected")

	if config.ENV == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db.AutoMigrate(&models.User{})

	route := gin.Default()
	route.Run(config.PORT)
}