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

func (api *API) getPrepaidCard(c *gin.Context) {
	userID, err := api.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	prepaidCards, err := api.prepaidRepo.GetPrepaidCardByUserID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	var response []dto.PrepaidCardResponse

	for i := range prepaidCards {
		response = append(response, dto.PrepaidCardResponse{
			ID:          prepaidCards[i].ID,
			Title:       prepaidCards[i].Title,
			ServiceID:   prepaidCards[i].ServiceID,
			StatusID:    prepaidCards[i].StatusID,
			DeadlineDay: prepaidCards[i].DeadlineDay,
			Amount:      prepaidCards[i].Amount,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (api *API) updatePrepaidCard(c *gin.Context) {
	var req dto.UpdatePrepaidRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	if err := api.prepaidRepo.UpdatePrepaidByID(req.ID, req.DeadlineDay, req.IdentityNumber, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SimpleResponse{Message: "Prepaid card updated successfully"})
}
