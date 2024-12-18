package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/sgitwhyd/music-catalogue/internal/configs"
	"github.com/sgitwhyd/music-catalogue/internal/handlers"
	"github.com/sgitwhyd/music-catalogue/internal/handlers/spotify"
	"github.com/sgitwhyd/music-catalogue/internal/models"
	spotifyModel "github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	"github.com/sgitwhyd/music-catalogue/internal/repositorys"
	spotifyRepo "github.com/sgitwhyd/music-catalogue/internal/repositorys/spotify"
	"github.com/sgitwhyd/music-catalogue/internal/services"
	spotifySvc "github.com/sgitwhyd/music-catalogue/internal/services/spotify"
	"github.com/sgitwhyd/music-catalogue/pkg/httpclient"
	"github.com/sgitwhyd/music-catalogue/pkg/internalsql"
)

func main() {
	config, err := configs.Init("./", "env", ".env")
	if err != nil {
		log.Error().Err(err).Msgf("load config error %v", err)
	}
	
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
	db.AutoMigrate(&spotifyModel.TrackActivity{})


	r := gin.Default()
	route :=r.Group("/api/v1")

	client := httpclient.NewClient(&http.Client{})
	if client == nil {
    log.Fatal().Msg("Failed to initialize http client")
	}

	log.Info().Msgf("client %v", client)

	// repositorys
	spotifyOutbond := spotifyRepo.NewSpotifyOutbond(config, client)
	userRepo := repositorys.NewUserRepo(db)
	spotifyRepository := spotifyRepo.NewSpotifyRepository(db)


	// services
	userService := services.NewUserService(userRepo, config)
	spotifyService := spotifySvc.NewSpotifyServie(spotifyOutbond, spotifyRepository)

	// // handlers
	userHandler := handlers.NewUserHandler(userService, route)
	spotifyHandler := spotify.NewSpotifyHandler(spotifyService, route)

	// // register route
	userHandler.RegisterRoute()
	spotifyHandler.RegisterRoute()

	r.Run(config.PORT)
}