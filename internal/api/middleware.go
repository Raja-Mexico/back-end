package api

import (
	"net/http"

	"github.com/Raja-Mexico/back-end/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString = tokenString[len("Bearer "):]

	responseCode, err := validateToken(tokenString)

	if err != nil {
		c.JSON(responseCode, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.Next()
}

func validateToken(tokenString string) (responseCode int, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized, err
		}
		return http.StatusBadRequest, err
	}

	if !token.Valid {
		return http.StatusUnauthorized, err
	}

	return http.StatusOK, nil
}
