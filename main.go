package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/piyushpatil22/dapter/dapter"
)

type User struct {
	ID        int    `json:"id" dapTableName:"clients2"  dapFieldAttrs:"PK , NotNull"`
	FirstName string `json:"firstName" dapFieldAttrs:"PK NotNull"`
	LastName  string `json:"lastName" dapFieldAttrs:"NotNull"`
}

func main() {

	connection := "host=localhost port=5432 user=postgres password=root dbname=dashapp sslmode=disable"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	//this is the impl example for the DAP
	dap := dapter.NewDAP(db)

	var user User
	err = dap.AutoMigrate(user)
	if err != nil {
		log.Fatal(err)
	}

	//a get query would look something like
	//dap.Get(&user, "id", 1)
}
