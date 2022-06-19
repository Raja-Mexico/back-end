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

type InstitutionListResponse struct {
	Data []institution `json:"data"`
}

type transaction struct {
	Description         string              `json:"description"`
	Date                string              `json:"date"`
	Direction           string              `json:"direction" binding:"required"`
	Amount              float64             `json:"amount" binding:"required"`
	TransactionCategory transactionCategory `json:"category" binding:"required"`
}

type institution struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type transactionCategory struct {
	Name string `json:"classification_group" binding:"required"`
}

type CategorizeTransactionResponse struct {
	TotalExpense        float64                       `json:"total_expense"`
	TransactionCategory []TransactionCategoryResponse `json:"expenses"`
	TopExpense          []TransactionCategoryResponse `json:"top_expenses"`
	Institution         []string                      `json:"institutions"`
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

type CreateTeamRequest struct {
	Name string `json:"family_name" binding:"required"`
}

type CreateTeamResponse struct {
	Message    string `json:"message"`
	FamilyCode string `json:"family_code"`
}

type JoinDetailTeamRequest struct {
	FamilyCode string `json:"family_code" binding:"required"`
}

type JoinTeamResponse struct {
	Success bool `json:"success"`
}

type DetailTeamResponse struct {
	FamilyCode    string               `json:"family_code"`
	FamilyBalance float64              `json:"family_balance"`
	Members       []TeamMemberResponse `json:"members"`
}

type TeamMemberResponse struct {
	Name     string  `json:"name"`
	Balance  float64 `json:"balance"`
	IsAdmin  bool    `json:"is_admin"`
	IsSender bool    `json:"is_sender"`
}

type TeamExpenseResponse struct {
	Expenses []Expense
}

type Expense struct {
	Spender  string  `json:"spender"`
	Desc     string  `json:"desc"`
	Date     string  `json:"date"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}
