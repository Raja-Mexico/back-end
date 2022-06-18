package api

import (
	"net/http"

	"github.com/Raja-Mexico/back-end/internal/dto"
	"github.com/gin-gonic/gin"
)

func (api *API) getUserInfo(c *gin.Context) {
	userID, err := api.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := api.userRepo.GetUserInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UserInfoResponse{
		Name:               user.Name,
		Balance:            user.Balance,
		VirtualAccountCode: user.NoVirtualAccount,
		IsInFamily:         user.IsInFamily,
	})
}
