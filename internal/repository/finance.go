package repository

import "database/sql"

type FinancialRepository struct {
	db *sql.DB
}

func NewFinancialRepository(db *sql.DB) *FinancialRepository {
	return &FinancialRepository{
		db: db,
	}
}

func (f *FinancialRepository) InsertUserFinanceAccount(userID, bankID int, accessToken string) error {
	statement := `INSERT INTO financial_account (user_id, bank_id, access_token) VALUES (?, ?, ?);`
	_, err := f.db.Exec(statement, userID, bankID, accessToken)
	return err
}

func (f *FinancialRepository) GetAccessTokenByUserID(userID int) ([]string, error) {
	statement := `SELECT access_token FROM financial_account WHERE user_id = ?;`

	rows, err := f.db.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accessTokens []string

	for rows.Next() {
		var accessToken string
		err := rows.Scan(&accessToken)
		if err != nil {
			return nil, err
		}

		accessTokens = append(accessTokens, accessToken)
	}

	return accessTokens, nil
}
