package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

type TeamRepository struct {
	db *sql.DB
}

type TeamMember struct {
	ID      int     `json:"id"`
	Name    string  `db:"name"`
	Balance float64 `db:"balance"`
	IsAdmin bool    `db:"is_admin"`
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (t *TeamRepository) CreateTeam(name string, userID int) (string, error) {
	teamID := uuid.New().String()[:8]

	for {
		isExist, err := t.CheckTeamExists(teamID)
		if err != nil {
			return "", err
		}

		if !isExist {
			break
		}

		teamID = uuid.New().String()[:8]
	}

	var balance float64
	statement := `SELECT balance FROM user_balance WHERE user_id = ?;`
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
		VALUES (?, ?, ?, ?);
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

func (t *TeamRepository) JoinTeam(teamID string, userID int) error {
	var userBalance float64
	statement := `SELECT balance FROM user_balance WHERE user_id = ?;`
	err := t.db.QueryRow(statement, userID).Scan(&userBalance)
	if err != nil {
		return err
	}

	statement = `DELETE FROM team WHERE creator_id = ?;`
	_, err = t.db.Exec(statement, userID)
	if err != nil {
		return err
	}

	statement = `
		UPDATE team SET balance = balance + ? WHERE id = ?;
	`
	_, err = t.db.Exec(statement, userBalance, teamID)
	if err != nil {
		return err
	}

	statement = `
		UPDATE membership SET team_id = ?, is_admin = ? WHERE user_id = ?;
	`
	_, err = t.db.Exec(statement, teamID, false, userID)
	if err != nil {
		return err
	}

	statement = `
		UPDATE user_balance SET team_id = ? WHERE user_id = ?;
	`
	_, err = t.db.Exec(statement, teamID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TeamRepository) CheckTeamExists(teamID string) (bool, error) {
	statement := `SELECT COUNT(*) FROM team WHERE id = ?;`

	var count int
	err := t.db.QueryRow(statement, teamID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (t *TeamRepository) GetMembers(teamID string) ([]TeamMember, error) {
	statement := `
	SELECT users.id, users.name, user_balance.balance, membership.is_admin
	FROM users
	INNER JOIN membership ON users.id = membership.user_id
	INNER JOIN user_balance ON user_balance.user_id = users.id
	WHERE membership.team_id = ?;
	`

	rows, err := t.db.Query(statement, teamID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	members := []TeamMember{}
	for rows.Next() {
		var member TeamMember
		err = rows.Scan(&member.ID, &member.Name, &member.Balance, &member.IsAdmin)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

func (t *TeamRepository) GetTeamByUserID(userID int) (string, error) {
	statement := `SELECT team_id FROM membership WHERE user_id = ?;`

	var teamID string
	err := t.db.QueryRow(statement, userID).Scan(&teamID)
	if err != nil {
		return "", err
	}

	return teamID, nil
}

func (t *TeamRepository) GetTeamName(teamID string) (string, error) {
	statement := `SELECT name FROM team WHERE id = ?;`

	var name string
	err := t.db.QueryRow(statement, teamID).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (t *TeamRepository) GetTeamBalance(teamID string) (float64, error) {
	statement := `SELECT balance FROM team WHERE id = ?;`

	var balance float64
	err := t.db.QueryRow(statement, teamID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
