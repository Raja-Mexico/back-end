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
	var req dto.JoinDetailTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	isTeamExist, err := api.teamRepo.CheckTeamExists(req.FamilyCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	if !isTeamExist {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "Team does not exist"})
		return
	}

	userID, err := api.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	if err := api.teamRepo.JoinTeam(req.FamilyCode, userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.JoinTeamResponse{Success: true})
}

func (api *API) getDetailTeam(c *gin.Context) {
	senderUserID, err := api.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	teamID, err := api.teamRepo.GetTeamByUserID(senderUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	teamName, err := api.teamRepo.GetTeamName(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	members, err := api.teamRepo.GetMembers(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	familyBalance, err := api.teamRepo.GetTeamBalance(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	var response dto.DetailTeamResponse
	response.FamilyName = teamName
	response.FamilyCode = teamID
	response.FamilyBalance = familyBalance
	for _, member := range members {
		response.Members = append(response.Members, dto.TeamMemberResponse{
			Name:     member.Name,
			Balance:  member.Balance,
			IsAdmin:  member.IsAdmin,
			IsSender: member.ID == senderUserID,
		})
	}

	c.JSON(http.StatusOK, response)
}
