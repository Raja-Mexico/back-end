package api

import (
	"net/http"

	"github.com/Raja-Mexico/back-end/internal/dto"
	"github.com/gin-gonic/gin"
)

func (api *API) createTeam(c *gin.Context) {
	var req dto.CreateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	userID, err := api.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	teamID, err := api.teamRepo.CreateTeam(req.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CreateTeamResponse{Message: "Team created successfully", FamilyCode: teamID})
}

func (api *API) joinTeam(c *gin.Context) {

}
