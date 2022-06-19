package api

import (
	"net/http"

	"github.com/Raja-Mexico/back-end/internal/dto"
	"github.com/gin-gonic/gin"
)

func (api *API) savePrepaidCard(c *gin.Context) {
	var req dto.CreatePrepaidRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	userID, err := api.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	teamID, err := api.teamRepo.GetTeamByUserID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	if err := api.prepaidRepo.InsertNewPrepaid(
		userID, req.ServiceID, teamID, req.DeadlineDay, req.IdentityNumber, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SimpleResponse{Message: "Prepaid card created successfully"})

}
