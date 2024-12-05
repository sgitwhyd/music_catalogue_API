package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/music-catalogue/internal/models"
	"github.com/sgitwhyd/music-catalogue/internal/services"
)

//go:generate mockgen -source=user_handler.go -destination=user_handler_mock_test.go -package=handlers
type UserService interface {
	services.UserService
}

type userHandler struct {
	userService services.UserService
	route *gin.RouterGroup
}

func NewUserHandler(userService services.UserService, route *gin.RouterGroup) *userHandler {
	return &userHandler{
		userService: userService,
		route: route,
	}
}

func (h *userHandler) SignUp(c *gin.Context) {
	var request models.SignUpRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Error().Err(err).Msgf("error hhhh %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}

	err = h.userService.Register(request)
	if err != nil {
		log.Error().Err(err).Msgf("error %v", err)
		if err.Error() == "username or email already registered" {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
				"errors": err.Error(),
			})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": "created",
	})
}

func (h *userHandler) RegisterRoute(){
	route := h.route.Group("/auth")

	route.POST("/signup", h.SignUp)
}