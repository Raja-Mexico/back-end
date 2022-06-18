package repository

import (
	"database/sql"
	"regexp"
	"unicode"

	"github.com/Raja-Mexico/back-end/internal/constant"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	EMAIL_REGEX = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
)

type UserInfo struct {
	Name             string  `db:"name"`
	Balance          float64 `db:"balance"`
	NoVirtualAccount string  `db:"no_virtual_account"`
	IsInFamily       bool
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) InsertNewUser(name, email, password string) (int, error) {
	isEmailValidate, err := u.isEmailValidated(email)

	if err != nil {
		return 0, err
	}

	if !isEmailValidate {
		return 0, constant.ErrEmailInvalid
	}

	isEmailExist, err := u.isEmailExist(email)
	if err != nil {
		return 0, err
	}

	if isEmailExist {
		return -1, constant.ErrEmailAlreadyExist
	}

	if !u.isPasswordValidated(password) {
		return 0, constant.ErrPasswordInvalid
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	statement := `INSERT INTO users (name, email, password) VALUES (?, ?, ?);`
	res, err := u.db.Exec(statement, name, email, hashedPassword)
	if err != nil {
		return 0, err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}

func (u *UserRepository) InjectUserNeedsAfterRegister(userID int) error {
	teamID := uuid.New().String()

	statement := `INSERT INTO team (id, creator_id) VALUES (?, ?);`
	_, err := u.db.Exec(statement, teamID, userID)
	if err != nil {
		return err
	}

	statement = `INSERT INTO membership (team_id, user_id) VALUES (?, ?);`
	_, err = u.db.Exec(statement, teamID, userID)
	if err != nil {
		return err
	}

	statement = `INSERT INTO user_balance (team_id, user_id, no_virtual_account) VALUES (?, ?, ?);`
	no_virtual_account := "130432138889902"
	_, err = u.db.Exec(statement, teamID, userID, no_virtual_account)
	return err
}

func (u *UserRepository) GetUserInfo(userID int) (UserInfo, error) {
	statement := `
	SELECT name, balance, no_virtual_account 
	FROM users 
	JOIN user_balance ON users.id = user_balance.user_id
	WHERE users.id = ?;`

	row := u.db.QueryRow(statement, userID)
	var userInfo UserInfo
	err := row.Scan(&userInfo.Name, &userInfo.Balance, &userInfo.NoVirtualAccount)
	if err != nil {
		return UserInfo{}, err
	}

	statement = `SELECT name FROM team WHERE creator_id = ?;`
	row = u.db.QueryRow(statement, userID)
	var teamName string
	err = row.Scan(&teamName)

	if err != nil {
		userInfo.IsInFamily = false
	} else {
		userInfo.IsInFamily = true
	}

	return userInfo, nil
}

func (u *UserRepository) CheckUserByEmailAndPassword(email, password string) (int, error) {
	isEmailExist, err := u.isEmailExist(email)
	if err != nil {
		return 0, err
	}

	if !isEmailExist {
		return -1, constant.ErrEmailNotFound
	}

	statement := `SELECT id, password FROM users WHERE email = ?;`

	var (
		userID         int
		hashedPassword string
	)

	err = u.db.QueryRow(statement, email).Scan(&userID, &hashedPassword)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return 0, constant.ErrPasswordIsWrong
	}

	return userID, nil
}

func (u *UserRepository) isPasswordValidated(password string) bool {
	if len(password) < 8 {
		return false
	}

	var isOneLetter, isOneNumber bool
	for _, c := range password {
		if unicode.IsLetter(c) {
			isOneLetter = true
		} else if unicode.IsNumber(c) {
			isOneNumber = true
		}
	}

	return isOneLetter && isOneNumber
}

func (u *UserRepository) isEmailValidated(email string) (bool, error) {
	regex, err := regexp.Compile(EMAIL_REGEX)
	if err != nil {
		return false, err
	}

	isValid := regex.MatchString(email)
	return isValid, nil
}

func (u *UserRepository) isEmailExist(email string) (bool, error) {
	statement := `SELECT COUNT(*) FROM users WHERE email = ?;`
	res := u.db.QueryRow(statement, email)

	var count int
	err := res.Scan(&count)
	if count > 0 {
		return true, err
	}
	return false, err
}
