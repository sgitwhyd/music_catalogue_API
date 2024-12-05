package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/sgitwhyd/music-catalogue/internal/configs"
	"github.com/sgitwhyd/music-catalogue/internal/handlers"
	"github.com/sgitwhyd/music-catalogue/internal/models"
	"github.com/sgitwhyd/music-catalogue/internal/repositorys"
	"github.com/sgitwhyd/music-catalogue/internal/services"
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

	// migrate db
	db.AutoMigrate(&models.User{})


	r := gin.Default()
	route :=r.Group("/api/v1")

	// repositorys
	userRepo := repositorys.NewUserRepo(db)

	// services
	userService := services.NewUserService(userRepo)

	// handlers
	userHandler := handlers.NewUserHandler(userService, route)

	// register route
	userHandler.RegisterRoute()

	r.Run(config.PORT)
}