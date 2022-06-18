package api

import (
	"github.com/Raja-Mexico/back-end/internal/repository"
	"github.com/gin-gonic/gin"
)

type API struct {
	userRepo      *repository.UserRepository
	financialRepo *repository.FinancialRepository
	router        *gin.Engine
}

func NewAPI(
	userRepo *repository.UserRepository,
	financialRepo *repository.FinancialRepository,
) *API {
	router := gin.Default()

	api := &API{
		router:        router,
		financialRepo: financialRepo,
		userRepo:      userRepo,
	}

	router.POST("/api/register", api.register)
	router.POST("/api/login", api.login)

	router.POST("/api/brick", api.postFinancialAccount)

	router.Use(
		authMiddleware,
	)
	{
		routerBrick := router.Group("/api/brick")
		{
			routerBrick.GET("/", api.getBrick)
			routerBrick.GET("/transaction", api.categorizeTransaction)
		}

	}

	return api
}

func (api *API) Handler() *gin.Engine {
	return api.router
}

func (api *API) Start() {
	api.Handler().Run()
}
