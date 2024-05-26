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

```
func main() {

	connection := "host=localhost port=5432 user=postgres password=root dbname=dbname sslmode=disable"
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
```
