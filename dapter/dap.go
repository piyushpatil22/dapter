package dapter

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// DAP is a Db wrapper to perform DB operations and get result in go  struct
// format. It uses the sql.DB package to perform database operations.
// It also uses the reflect package to parse the struct fields and validate
// the `dap` tag.
// The `dapTableName` tag is used to specify the table name for that struct/entity.
// and the `dapFieldAttrs`
// The `dapFieldAttrs` tag is used to specify the field attributes.
type DAP struct {
	db *sql.DB
}

// Creates a new DAP instance.
// It takes a pointer to a sql.DB instance as an argument and returns a
// pointer to a DAP instance.
func NewDAP(db *sql.DB) *DAP {
	return &DAP{
		db: db,
	}
}

// Closes the database connection via DAP
func (d *DAP) Close() error {
	return d.db.Close()
}

// Parses the struct that is passed in and validates the `dap` tag in the
// struct fields.
// The struct will be considered a DB entity if the `dapTableName` tag is
// found in the ID column
// Returns an error if the `dap` tag is not found or the `dapTableName`
// tag is not found.
func (d *DAP) AutoMigrate(object interface{}) error {
	objectValue := reflect.ValueOf(object)
	objectType := objectValue.Type()

	//TODO might change this from ID to the first field in that struct
	field, _ := objectType.FieldByName("ID")
	if tableName := field.Tag.Get("dapTableName"); tableName != "" {
		pkFieldName, pkErr := CheckPKDefined(objectValue)
		if pkErr != nil {
			return pkErr
		}
		//TODO: add this PK attr to that field
		log.Printf("Field with Primary Key: %v ", pkFieldName)
		var exists bool
		err := d.db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = $1)", tableName).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			log.Printf("table %v exists", tableName)
			//TODO: do you want to do migration here?
		} else {
			log.Printf("Table %s does not exist\n", tableName)
			query, err := createTableCreationQuery(objectValue, tableName, pkFieldName)
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

// Returns true if only one PK is defined in the struct,
// false if more than one Pk is defined or no PK is defined
func CheckPKDefined(object reflect.Value) (string, error) {
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
		return fieldName, nil
	}
	if pkOccurences > 1 {
		return "", ErrMultiplePKFieldFound
	}
	return "", fmt.Errorf("Error while determining PK")
}

// Creates a query to create a table based on the struct that is passed in.
// Returns an error if the struct is empty/invalid.
func createTableCreationQuery(object reflect.Value, tableName string, pkFieldName string) (string, error) {
	createTableQueryPrefix := "CREATE TABLE " + tableName + " ("
	querySuffix := ");"
	query := ""
	fieldNums := object.NumField()
	if fieldNums == 0 {
		return createTableQueryPrefix + querySuffix, nil
	} else {
		for i := 0; i < fieldNums; i++ {
			fieldType := object.Type().Field(i).Type
			fieldName := object.Type().Field(i).Name
			PK := ""
			if fieldName == pkFieldName {
				PK = " PRIMARY KEY"
			}
			if postgresType, ok := goTypeToPostgresType[fieldType.Name()]; ok {
				query += object.Type().Field(i).Name + " " + postgresType + PK
			}
			if i != fieldNums-1 {
				query += ", "
			}
			log.Println(query)
		}
	}
	return createTableQueryPrefix + query + querySuffix, nil
}
