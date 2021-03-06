package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Raja-Mexico/back-end/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Id int
	jwt.StandardClaims
}

var jwtKey = []byte("SECRET_KEY")

func (api *API) register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	userID, err := api.userRepo.InsertNewUser(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	tokenString, err := api.generateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	err = api.userRepo.InjectUserNeedsAfterRegister(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "User created successfully",
		Token:   tokenString,
	})
}

func (api *API) login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	userID, err := api.userRepo.CheckUserByEmailAndPassword(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	tokenString, err := api.generateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "User logged in successfully",
		Token:   tokenString,
	})
}

func (api *API) generateJWT(userId int) (string, error) {
	expTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Id: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func (api *API) getUserIDFromToken(c *gin.Context) (int, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return 0, fmt.Errorf("token not found")
	}

	claims := &Claims{}
	tokenString = tokenString[len("Bearer "):]
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, err
	}

	return claims.Id, nil
}
