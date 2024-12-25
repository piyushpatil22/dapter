package main

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/piyushpatil22/dapter/dap"
	"github.com/piyushpatil22/dapter/dap/builder"
	"github.com/piyushpatil22/dapter/dap/filter"
	"github.com/piyushpatil22/dapter/dap/store"
	"github.com/piyushpatil22/dapter/log"
)

type Base struct {
	ID        int64     `json:"id"`
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

	store := store.NewStore(db)
	defer store.Close()

	user := User{
		Username:       "sam",
		Password:       "notsecure",
		Gender:         "male",
		AccountBalance: 555,
		IsAdmin:        true,
	}

	//inserting entity
	err = store.Insert(user)
	if err != nil {
		log.Log.Err(err).Msg("Error inserting user")
	}

	//conditioning logic type 1
	filter := filter.Filter{
		Field: "username",
		Value: "sam",
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
	log.Log.Info().Interface("list", list).Msg("User List fetched via simple conditions")

	//conditioning logic type 2
	newFilters := []builder.Condition{
		builder.QueryCondition{Field: "account_balance", Operator: ">", Value: "100"},
		builder.QueryCondition{Field: "id", Operator: "=", Value: "69"},
	}

	//TODO need to figure out this part (either rewrite the builder package or get everything together)
	var userList User
	args, query := builder.Select(&userList, User{}).WHERE(newFilters).GetQueryArgs()
	log.Log.Info().Str("query", query).Interface("args", args).Msg("Query and Args")
	err = store.ExecuteWithConditions(&userList, query, args)
	if err != nil {
		if err == dap.ErrNoRowsFound {
			log.Log.Info().Msg("No rows found")
			return
		}
		log.Log.Err(err).Msg("Error getting user")
	}
	log.Log.Info().Interface("userList", userList).Msg("User List fetched via complex conditions")
}
