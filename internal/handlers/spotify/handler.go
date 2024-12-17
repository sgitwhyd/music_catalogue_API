package spotify

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/music-catalogue/internal/services/spotify"
)

type handler struct {
	service spotify.SpotifyService
	*gin.RouterGroup
}

func NewSpotifyHandler(service spotify.SpotifyService, route *gin.RouterGroup) *handler {
	return &handler{
		service: service,
		RouterGroup: route,
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

	response, err := h.service.Search(ctx, query, pageSize, pageIndex)
	if err != nil {
		log.Error().Err(err).Msg("error count")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,response)
}


func (h *handler) RegisterRoute(){
	route := h.RouterGroup.Group("/spotify")
	
	route.GET("/search", h.Search)

}