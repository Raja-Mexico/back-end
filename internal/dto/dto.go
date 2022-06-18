package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

type CredentialsFinancialAccount struct {
	BankID      string `json:"bankId" binding:"required"`
	AccessToken string `json:"accessToken" binding:"required"`
	UserID      string `json:"userId" binding:"required"`
}

type DataTransaction struct {
	Data []transaction `json:"data"`
}

type transaction struct {
	Direction           string              `json:"direction" binding:"required"`
	Amount              float64             `json:"amount" binding:"required"`
	TransactionCategory transactionCategory `json:"category" binding:"required"`
}

type transactionCategory struct {
	Name string `json:"classification_group" binding:"required"`
}

type CategorizeTransactionResponse struct {
	TotalExpense        float64                       `json:"total_expense"`
	TransactionCategory []TransactionCategoryResponse `json:"expenses"`
	TopExpense          []TransactionCategoryResponse `json:"top_expenses"`
}

type TransactionCategoryResponse struct {
	Name         string  `json:"name"`
	TotalExpense float64 `json:"total_expense"`
	Percentage   string  `json:"percentage"`
}

type UserInfoResponse struct {
	Name               string  `json:"name"`
	Balance            float64 `json:"balance"`
	VirtualAccountCode string  `json:"virtual_account_code"`
	IsInFamily         bool    `json:"is_in_family"`
}
