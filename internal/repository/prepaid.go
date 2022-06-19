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

type PrepaidCard struct {
	ID          int     `db:"id"`
	Title       string  `db:"title"`
	ServiceID   int     `db:"service_id"`
	StatusID    int     `db:"status_id"`
	DeadlineDay int     `db:"deadline_day"`
	Amount      float64 `db:"amount"`
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

func (p *PrepaidRepository) GetPrepaidCardByUserID(userID int) ([]PrepaidCard, error) {
	statement := `
		SELECT
		prepaid_card.id, service.name as title, service.id as service_id,
		status.id as status_id, prepaid_card.deadline_day, prepaid_card.nominal as amount
		FROM prepaid_card
		INNER JOIN status ON prepaid_card.status_id = status.id
		INNER JOIN Service ON prepaid_card.service_id = service.id
		WHERE prepaid_card.creator_id = ?;
	`
	rows, err := p.db.Query(statement, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var prepaidCards []PrepaidCard

	for rows.Next() {
		var prepaidCard PrepaidCard

		err := rows.Scan(
			&prepaidCard.ID,
			&prepaidCard.Title,
			&prepaidCard.ServiceID,
			&prepaidCard.StatusID,
			&prepaidCard.DeadlineDay,
			&prepaidCard.Amount,
		)

		if err != nil {
			return nil, err
		}

		prepaidCards = append(prepaidCards, prepaidCard)
	}

	return prepaidCards, nil

}

func (p *PrepaidRepository) GetPrepaidCardByID(cardID int) (PrepaidCard, error) {
	statement := `
		SELECT
		prepaid_card.id, service.name as title, service.id as service_id,
		status.id as status_id, prepaid_card.deadline_day, prepaid_card.nominal as amount
		FROM prepaid_card
		INNER JOIN status ON prepaid_card.status_id = status.id
		INNER JOIN Service ON prepaid_card.service_id = service.id
		WHERE prepaid_card.id = ?;
	`
	row := p.db.QueryRow(statement, cardID)

	var prepaidCard PrepaidCard

	err := row.Scan(
		&prepaidCard.ID,
		&prepaidCard.Title,
		&prepaidCard.ServiceID,
		&prepaidCard.StatusID,
		&prepaidCard.DeadlineDay,
		&prepaidCard.Amount,
	)

	if err != nil {
		return prepaidCard, err
	}

	return prepaidCard, nil
}

func (p *PrepaidRepository) UpdatePrepaidByID(
	id int, deadlineDay int, destinationNumber string, nominal float64) error {
	statement := `
		UPDATE prepaid_card SET deadline_day = ?, destination_number = ?, nominal = ?
		WHERE id = ?;
	`
	_, err := p.db.Exec(statement, deadlineDay, destinationNumber, nominal, id)

	if err != nil {
		return err
	}

	return nil
}

func (p *PrepaidRepository) PayPrepaidByID(id, userID int, teamID, destinationNumber string, amount float64) error {
	statement := `UPDATE prepaid_card 
	SET status_id = 2, nominal = ?, destination_number = ?
	WHERE id = ?;`
	_, err := p.db.Exec(statement, amount, destinationNumber, id)
	if err != nil {
		return err
	}

	statement = `UPDATE user_balance SET balance = balance - ? WHERE id = ?;`
	_, err = p.db.Exec(statement, amount, userID)
	if err != nil {
		return err
	}

	statement = `UPDATE team SET prepaid_balance = prepaid_balance + ?, balance = balance - ? WHERE id = ?;`
	_, err = p.db.Exec(statement, amount, amount, teamID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PrepaidRepository) RequestPrepaidCardPay(userID, prepaidCardID int, payRequest float64) error {
	statement := `INSERT INTO manage_prepaid_card (user_id, prepaid_card_id, pay_request) VALUES (?, ?, ?);`
	_, err := p.db.Exec(statement, userID, prepaidCardID, payRequest)
	if err != nil {
		return err
	}

	statement = `UPDATE prepaid_card SET status_id = 3 WHERE id = ?;`
	_, err = p.db.Exec(statement, prepaidCardID)
	if err != nil {
		return err
	}

	return nil
}
