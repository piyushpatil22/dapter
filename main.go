package main

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/piyushpatil22/dapter/dap"
	"github.com/piyushpatil22/dapter/dap/filter"
	"github.com/piyushpatil22/dapter/log"
)

type Base struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	Base
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	Email          string  `json:"email"`
	Phone          int64   `json:"phone"`
	AccountBalance float64 `json:"account_balance"`
	Gender         string  `json:"gender"`
	DOB            string  `json:"dob"`
	IsActivated    bool    `json:"is_activated"`
	IsAdmin        bool    `json:"is_admin"`
}

func main() {
	_ = log.Log
	connection := "host=localhost port=5432 user=postgres password=root dbname=dapter_test sslmode=disable"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Log.Err(err).Msg("Error connecting to database")
	}

	store := dap.NewStore(db)
	defer store.Close()

	user := User{
		Username: "jack",
		Password: "lantern",
	}

	err = store.Insert(user)
	if err != nil {
		log.Log.Err(err).Msg("Error inserting user")
	}

	filter := filter.Filter{
		Field: "username",
		Value: "jack",
	}
	var list []User
	err = store.GetByFilter(&list, filter, User{})
	if err != nil {
		if err == dap.ErrNoRowsFound {
			log.Log.Info().Msg("No rows found")
			return
		}
		log.Log.Err(err).Msg("Error getting user")
	}
	for _, u := range list {
		log.Log.Info().Interface("user", u).Msg("User")
	}
}
