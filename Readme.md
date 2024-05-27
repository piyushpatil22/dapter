# Dapter

## Generic DB Wrapper for Go Projects

## Goal

The goal of this project is to create a generic wrapper for database operations in Go projects. This wrapper should simplify database interactions and promote code reuse across different projects.

## Potential Features

-   Perform CRUD (Create, Read, Update, Delete) operations on database entities.
-   Flexibility to work with different database systems.
-   Simple API for querying and retrieving data.
-   Support for transactions.

## An Example

```go
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/piyushpatil22/dapter/dapter"
)

type User struct {
	ID        int    `json:"id" dapTableName:"user"  dapFieldAttrs:"PK , NotNull"`
	FirstName string `json:"firstName" dapFieldAttrs:"NotNull"`
	LastName  string `json:"lastName" dapFieldAttrs:"NotNull"`
}

func main() {

	connection := "host=localhost port=5432 user=postgres password=root dbname=dbname sslmode=disable"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	//this is the impl example for the DAP
	dap := dapter.NewDAP(db)

	var user User

	// this will create a user table by parsing the tags
	// from the user struct. Like primary key, not null
	// DB type will be determined from the go type.
	err = dap.AutoMigrate(user)
	if err != nil {
		log.Fatal(err)
	}

	//a get query would look something like
	//dap.Get(&user, "id", 1)
}
```
