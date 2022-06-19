package api

import (
	"github.com/Raja-Mexico/back-end/internal/repository"
	"github.com/gin-gonic/gin"
)

type API struct {
	userRepo      *repository.UserRepository
	financialRepo *repository.FinancialRepository
	teamRepo      *repository.TeamRepository
	router        *gin.Engine
}

func NewAPI(
	userRepo *repository.UserRepository,
	financialRepo *repository.FinancialRepository,
	teamRepo *repository.TeamRepository,
) *API {
	router := gin.Default()

	api := &API{
		router:        router,
		financialRepo: financialRepo,
		userRepo:      userRepo,
		teamRepo:      teamRepo,
	}

	router.POST("/api/register", api.register)
	router.POST("/api/login", api.login)

	router.POST("/api/brick", api.postFinancialAccount)

	router.Use(
		authMiddleware,
	)
	{
		router.GET("/api/user-info", api.getUserInfo)

		routerTeam := router.Group("/api/team")
		{
			routerTeam.GET("/:id", api.getDetailTeam)
			routerTeam.POST("/", api.createTeam)
			routerTeam.POST("/join", api.joinTeam)
		}

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
