package spotify

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/music-catalogue/internal/middleware"
	"github.com/sgitwhyd/music-catalogue/internal/models/spotify"
	spotifyService "github.com/sgitwhyd/music-catalogue/internal/services/spotify"
)

type handler struct {
	service spotifyService.SpotifyService
	route *gin.RouterGroup
}

func NewSpotifyHandler(service spotifyService.SpotifyService, route *gin.RouterGroup) *handler {
	return &handler{
		service: service,
		route: route,
	}
}

func (h *handler) Search(c *gin.Context){
	ctx := c.Request.Context()

	query := c.Query("query")
	pageIndexStr := c.Query("pageIndex")
	pageSizeStr := c.Query("pageSize")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		pageIndex = 1
	}

	userID := c.GetUint("userID")

	response, err := h.service.Search(ctx, query, pageSize, pageIndex, userID)
	if err != nil {
		log.Error().Err(err).Msg("error count")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,response)
}

func (h *handler) UpsertActivity(c *gin.Context){
	ctx := c.Request.Context()

	var request spotify.TrackActivityRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetUint("userID")
	err = h.service.UpSertActivity(ctx, userID, request )
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{

			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "created",
	})
}


func (h *handler) RegisterRoute(){
	route := h.route.Group("/spotify")
	route.Use(middleware.AuthMiddleware())
	
	route.GET("/search", h.Search)
	route.POST("/activity", h.UpsertActivity)
	

}