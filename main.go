package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type User struct {
	ID        int    `json:"id"  dapFieldAttrs:"NotNull"`
	FirstName string `json:"firstName" dapFieldAttrs:" PK, NotNull"`
	LastName  string `json:"lastName" dapFieldAttrs:"NotNull"`
}

type UserDetails struct {
	ID      int    `json:"id" dapTableName:"clients_details" dapFieldAttrs:"PK, NotNull"`
	UserID  int    `json:"userId" dapFieldAttrs:"FK, NotNull"`
	Email   string `json:"email" dapFieldAttrs:"NotNull"`
	Phone   string `json:"phone" dapFieldAttrs:"NotNull"`
	Address string `json:"address" dapFieldAttrs:"NotNull"`
}

type UserAddress struct {
	ID      int    `json:"id" dapTableName:"clients_address" dapFieldAttrs:"PK, NotNull"`
	UserID  int    `json:"userId" dapFieldAttrs:"FK, NotNull"`
	Address string `json:"address" dapFieldAttrs:"NotNull"`
}

type DAP struct {
	db *sql.DB
}

var goTypeToPostgresType = map[string]string{
	"int":       "INTEGER",
	"int8":      "INTEGER",
	"int16":     "INTEGER",
	"int32":     "INTEGER",
	"int64":     "BIGINT",
	"uint":      "INTEGER",
	"uint8":     "INTEGER",
	"uint16":    "INTEGER",
	"uint32":    "INTEGER",
	"uint64":    "BIGINT",
	"float32":   "REAL",
	"float64":   "DOUBLE PRECISION",
	"bool":      "BOOLEAN",
	"string":    "VARCHAR",
	"time.Time": "TIMESTAMP",
	"[]byte":    "BYTEA",
}

const (
	NOT_NULL = "NOT NULL"
	VARCHAR  = "VARCHAR NOT NULL"
)

func main() {
	// create DB connection for postgres
	connection := "host=localhost port=5432 user=postgres password=root dbname=dashapp sslmode=disable"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	yoDap := DAP{db: db}

	var user User
	err = yoDap.Create(user)
	if err != nil {
		log.Fatal(err)
	}

	//create table

}

func (d *DAP) Create(object interface{}) error {
	objectValue := reflect.ValueOf(object)
	objectType := objectValue.Type()

	field, _ := objectType.FieldByName("ID")
	if tableName := field.Tag.Get("dapTableName"); tableName != "" {
		fieldName, ok := CheckPKDefined(objectValue)
		if !ok {
			return fmt.Errorf("error while parsing PK")
		}
		log.Printf("Field with Primary Key: %v ", fieldName)
		PK := field.Tag.Get("dapFieldAttrs")
		if PK != "" {
			if strings.Contains(PK, "PK") {
				log.Println("This table has a PK defined")
			}
		}
		log.Println(tableName)
		// Query to check table existence
		var exists bool
		err := d.db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = $1)", tableName).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			log.Printf("Table %s exists\n", tableName)
			// Query to get table schema
			// query := `SELECT column_name, data_type
			// 		FROM information_schema.columns
			// 		WHERE table_name = 'clients'`

			// // Execute the query
			// rows, err := d.db.Query(query)
			// if err != nil {
			// 	fmt.Println("Error executing query:", err)
			// 	return err
			// }
			// defer rows.Close()

			// // Iterate over the rows and print column names and types
			// for rows.Next() {
			// 	var columnName, dataType string
			// 	if err := rows.Scan(&columnName, &dataType); err != nil {
			// 		fmt.Println("Error scanning row:", err)
			// 		return err
			// 	}
			// 	fmt.Printf("Column Name: %s, Data Type: %s\n", columnName, dataType)
			// }
			// if err := rows.Err(); err != nil {
			// 	fmt.Println("Error iterating over rows:", err)
			// 	return err
			// }
			// log.Println("Table schema:")
		} else {
			log.Printf("Table %s does not exist\n", tableName)
			query, err := createTableCreationQuery(objectValue, tableName)
			if err != nil {
				return err
			}
			log.Println(query)

		}
	} else {
		return ErrDapTableTagNotFound
	}
	return nil
}

func createTableCreationQuery(object reflect.Value, tableName string) (string, error) {
	queryPrefix := "CREATE TABLE " + tableName + " ("
	querySuffix := ");"
	query := ""
	fieldNums := object.NumField()
	if fieldNums == 0 {
		log.Println("No fields found, simply creating a emprty table")
		return queryPrefix + querySuffix, nil
	} else {
		log.Println("Struct has fields: " + strconv.FormatInt(int64(fieldNums), 10))
		for i := 0; i < fieldNums; i++ {
			//determine the type of the field
			fieldType := object.Type().Field(i).Type
			if postgresType, ok := goTypeToPostgresType[fieldType.Name()]; ok {
				query += object.Type().Field(i).Name + " " + postgresType + " " + NOT_NULL
				// log.Printf("go struct value: %s to Postgres type %s", object.Type().Field(i).Name, postgresType)
			}
			if i != fieldNums-1 {
				query += ", "
			}
		}
	}
	return queryPrefix + query + querySuffix, nil
}

func CheckPKDefined(object reflect.Value) (string, bool) {
	var pkOccurences = 0
	var fieldName = ""
	fieldNums := object.NumField()
	for i := 0; i < fieldNums; i++ {
		field := object.Type().Field(i)
		if tableName := field.Tag.Get("dapFieldAttrs"); tableName != "" {
			if strings.Contains(tableName, "PK") {
				fieldName = field.Name
				pkOccurences++
			}
		}
	}
	if pkOccurences == 1 {
		return fieldName, true
	}
	if pkOccurences > 1 {
		return "", false
	}
	return "", false
}
