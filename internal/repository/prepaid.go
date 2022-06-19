package repository

import "database/sql"

type PrepaidRepository struct {
	db *sql.DB
}

func NewPrepaidRepository(db *sql.DB) *PrepaidRepository {
	return &PrepaidRepository{
		db: db,
	}
}

func (p *PrepaidRepository) InsertNewPrepaid(
	creatorID int, serviceID int, teamID string, deadlineDay int,
	destinationNumber string, nominal float64) error {
	statement := `
		INSERT INTO prepaid_card (creator_id, service_id, team_id, deadline_day, destination_number, nominal) VALUES 
		(?, ?, ?, ?, ?, ?);
	`
	_, err := p.db.Exec(statement, creatorID, serviceID, teamID, deadlineDay, destinationNumber, nominal)

	if err != nil {
		return err
	}

	return nil
}
