package main

import (
	"database/sql"

	"github.com/Raja-Mexico/back-end/internal/api"
	"github.com/Raja-Mexico/back-end/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "raja-mexico.db")
	if err != nil {
		panic(err)
	}

	userRepository := repository.NewUserRepository(db)
	financialRepository := repository.NewFinancialRepository(db)
	teamRepository := repository.NewTeamRepository(db)

	mainAPI := api.NewAPI(userRepository, financialRepository, teamRepository)
	mainAPI.Start()
}
