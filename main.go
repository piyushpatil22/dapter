package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/piyushpatil22/dapter/dapter"
)

type User struct {
	ID        int    `json:"id" dapTableName:"user"  dapFieldAttrs:"PK"`
	FirstName string `json:"firstName" dapFieldAttrs:"NOT NULL"`
	LastName  string `json:"lastName" dapFieldAttrs:"NOT NULL"`
}

type UserNew struct {
	ID            int       `json:"id" dapTableName:"usersnew" dapFieldAttrs:"PK"`
	FirstName     string    `json:"first_name" dapFieldAttrs:"NOT NULL"`
	LastName      string    `json:"last_name" dapFieldAttrs:"NOT NULL"`
	Email         string    `json:"email" dapFieldAttrs:"NOT NULL"`
	Password      string    `json:"-" dapFieldAttrs:"NOT NULL"`
	IsActive      bool      `json:"is_active"`
	CreatedOn     time.Time `json:"created_at"`
	LastLogin     time.Time `json:"last_login,omitempty"`
	Role          string    `json:"role,omitempty" dapFieldAttrs:"NOT NULL"`
	Birthday      time.Time `json:"birthday,omitempty"`
	DeactivatedOn time.Time `json:"deactivated_on,omitempty"`
}

func main() {

	connection := "host=localhost port=5432 user=postgres password=root dbname=dashapp sslmode=disable"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	//this is the impl example for the DAP
	dap := dapter.NewDAP(db)

	// var user User
	var userNew UserNew

	// err = dap.AutoMigrate(user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = dap.AutoMigrate(userNew)
	if err != nil {
		log.Fatal(err)
	}

	//a get query would look something like
	//dap.Get(&user, "id", 1)
}
