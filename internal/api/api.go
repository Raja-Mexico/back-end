package api

import (
	"github.com/Raja-Mexico/back-end/internal/repository"
	"github.com/gin-gonic/gin"
)

type API struct {
	userRepo *repository.UserRepository
	router   *gin.Engine
}

func NewAPI(
	userRepo *repository.UserRepository,
) *API {
	router := gin.Default()

	api := &API{
		router:   router,
		userRepo: userRepo,
	}

	router.POST("/api/register", api.register)
	router.POST("/api/login", api.login)

	return api
}

func (api *API) Handler() *gin.Engine {
	return api.router
}

func (api *API) Start() {
	api.Handler().Run()
}
