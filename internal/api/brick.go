package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/Raja-Mexico/back-end/internal/dto"
	"github.com/gin-gonic/gin"
)

var (
	accessToken = "public-sandbox-b7ed9e30-8109-4c4b-a895-8596bec10192"
	redirectUrl = "https://ce49-180-252-172-19.ap.ngrok.io/api/brick"
)

func (api *API) getBrick(c *gin.Context) {
	userID, err := api.getUserIDFromToken(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	launchBrickUrl := fmt.Sprintf("https://cdn.onebrick.io/sandbox-widget/v1/?accessToken=%s&redirect_url=%s&user_id=%d", accessToken, redirectUrl, userID)

	c.String(http.StatusOK, launchBrickUrl)
}

func (api *API) postFinancialAccount(c *gin.Context) {
	var req []dto.CredentialsFinancialAccount

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	go func() {
		for _, v := range req {
			userID, err := strconv.Atoi(v.UserID)
			if err != nil {
				log.Println(err)
				return
			}

			bankID, err := strconv.Atoi(v.BankID)
			if err != nil {
				log.Println(err)
				return
			}

			if err := api.financialRepo.InsertUserFinanceAccount(userID, bankID, v.AccessToken); err != nil {
				log.Println(err)
				return
			}
		}
	}()

	c.String(http.StatusOK, "OK")
}

func (api *API) categorizeTransaction(c *gin.Context) {
	userID, err := api.getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	banksID, err := api.financialRepo.GetUserFinanceBank(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	institutions, err := api.processInstitutionList()
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	institutionsMap := make(map[int]string)
	for _, institution := range institutions.Data {
		institutionsMap[institution.Id] = institution.Name
	}

	userBanks := make([]string, 0, len(banksID))
	for _, bankID := range banksID {
		userBanks = append(userBanks, institutionsMap[bankID])
	}

	accessTokens, err := api.financialRepo.GetAccessTokenByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	transactionCategory := make(map[string]float64)

	now := time.Now().Format("2006-01-02")
	yearAgo := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	url := fmt.Sprintf("https://sandbox.onebrick.io/v1/transaction/list?from=%s&to=%s", yearAgo, now)

	for _, v := range accessTokens {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
			return
		}

		req.Header.Set("Authorization", "Bearer "+v)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Println(resp.StatusCode)
			continue
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err.Error())
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
			return
		}

		var responseObject dto.DataTransaction
		json.Unmarshal(bodyBytes, &responseObject)

		for _, data := range responseObject.Data {
			if data.Direction == "out" {
				transactionCategory[data.TransactionCategory.Name] += data.Amount
			}
		}

	}

	top_expense, total := processData(transactionCategory)

	var response dto.CategorizeTransactionResponse
	response.Institution = userBanks
	response.TotalExpense = total
	for key, value := range transactionCategory {
		response.TransactionCategory = append(response.TransactionCategory, dto.TransactionCategoryResponse{
			Name:         key,
			TotalExpense: value,
			Percentage:   fmt.Sprintf("%.2f", (value/total)*100),
		})
	}
	for _, v := range top_expense {
		response.TopExpense = append(response.TopExpense, dto.TransactionCategoryResponse{
			Name:         v,
			TotalExpense: transactionCategory[v],
			Percentage:   fmt.Sprintf("%.2f", (transactionCategory[v]/total)*100),
		})
	}

	c.JSON(http.StatusOK, response)
}

func (api *API) processInstitutionList() (dto.InstitutionListResponse, error) {
	var institutions dto.InstitutionListResponse

	url := "https://sandbox.onebrick.io/v1/institution/list"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return institutions, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return institutions, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		return institutions, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err.Error())
		return institutions, err
	}

	json.Unmarshal(bodyBytes, &institutions)

	return institutions, nil
}

func processData(data map[string]float64) ([]string, float64) {
	total := 0.0
	keys := make([]string, 0, len(data))

	for k := range data {
		total += data[k]
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return data[keys[i]] > data[keys[j]]
	})

	if len(keys) > 3 {
		keys = keys[:3]
	}

	return keys, total
}
