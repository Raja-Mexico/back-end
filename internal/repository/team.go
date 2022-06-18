package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (t *TeamRepository) CreateTeam(name string, userID int) (string, error) {
	teamID := uuid.New().String()

	var balance float64
	statement := `SELECT balance FROM user_balance WHERE user_id = ?`
	err := t.db.QueryRow(statement, userID).Scan(&balance)
	if err != nil {
		return "", err
	}

	statement = `DELETE FROM team WHERE creator_id = ?;`
	_, err = t.db.Exec(statement, userID)
	if err != nil {
		return "", err
	}

	statement = `
		INSERT INTO team (id, name, creator_id, balance)
		VALUES (?, ?, ?, ?)
	`
	_, err = t.db.Exec(statement, teamID, name, userID, balance)
	if err != nil {
		return "", err
	}

	statement = `
		UPDATE membership SET team_id = ?, is_admin = ? WHERE user_id = ?;
	`
	_, err = t.db.Exec(statement, teamID, true, userID)
	if err != nil {
		return "", err
	}

	statement = `
		UPDATE user_balance SET team_id = ? WHERE user_id = ?;
	`
	_, err = t.db.Exec(statement, teamID, userID)
	if err != nil {
		return "", err
	}

	return teamID, nil
}
