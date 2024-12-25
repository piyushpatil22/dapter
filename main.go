package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/piyushpatil22/dapter/dap"
	"github.com/piyushpatil22/dapter/dap/filter"
	"github.com/piyushpatil22/dapter/log"
)

type Base struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type Instrument struct {
	Base
	Token            string `json:"token"`
	Symbol           string `json:"symbol"`
	Name             string `json:"name"`
	Expiry           string `json:"expiry"`
	StrikePrice      string `json:"strike_price"`
	LotSize          string `json:"lot_size"`
	InstrumentType   string `json:"instrument_type"`
	ExchaangeSegment string `json:"exch_seg"`
	TickSize         string `json:"tick_size"`
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
		Username:       "jack ",
		Password:       "jack123",
		Gender:         "male",
		AccountBalance: 1000,
	}

	err = store.Insert(user)
	if err != nil {
		log.Log.Err(err).Msg("Error inserting user")
	}

	filter := filter.Filter{
		Field: "gender",
		Value: "female",
	}
	userList := []User{}
	err = store.GetByFilter(&userList, filter)
	if err != nil {
		log.Log.Err(err).Msg("Error getting user")
	}

	// err = store.Update(user)
	// if err != nil {
	// 	log.Log.Err(err).Msg("Error updating user")
	// }

}
