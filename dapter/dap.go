package dapter

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// Interface for the DAP struct
// AutoMigrate - Auto migrate the struct to the database
// Create - Create a new record in the database
// Update - Update an existing record in the database
// Delete - Delete a record from the database
// Get - Get a record from the database
// GetAll - Get all records from the database
// Close - Close the database connection
type DAPInterface interface {
	AutoMigrate(object interface{}) error
	Create(object interface{}) error
	Update(object interface{}) error
	Delete(object interface{}) error
	Get(object interface{}, id string) error
	GetAll(object interface{}) error
	Close() error
}

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
		//get primary key and fail fast
		pkFieldName, pkErr := determinePrimaryKey(objectValue)
		if pkErr != nil {
			return pkErr
		}

		//get not null fields
		notNullFields, notNullErr := determineNotNullFields(objectValue)
		if notNullErr != nil {
			return notNullErr
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
			query, err := createTableCreationQuery(objectValue, tableName, pkFieldName, notNullFields)
			queryV2, errv2 := createTableCreationQueryV2(objectValue, tableName, pkFieldName, notNullFields)
			if errv2 != nil {
				return errv2
			}
			log.Println(queryV2)

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

// Creates a query to create a table based on the struct that is passed in.
// Returns an error if the struct is empty/invalid.
func createTableCreationQuery(object reflect.Value, tableName string, pkFieldName string, notNullFields []string) (string, error) {
	createTableQueryPrefix := "CREATE TABLE " + tableName + " ("
	querySuffix := ");"
	query := ""
	fieldNums := object.NumField()

	notNullMap := make(map[string]bool)
	for _, field := range notNullFields {
		notNullMap[field] = true
	}

	if fieldNums == 0 {
		return createTableQueryPrefix + querySuffix, nil
	} else {
		for i := 0; i < fieldNums; i++ {
			fieldType := object.Type().Field(i).Type
			fieldName := object.Type().Field(i).Name

			if fieldName == pkFieldName {
				query += fieldName + " SERIAL PRIMARY KEY"
			} else {
				if postgresType, ok := goTypeToPostgresType[fieldType.Name()]; ok {
					query += object.Type().Field(i).Name + " " + postgresType
				}
				if notNullMap[fieldName] {
					query += " NOT NULL"
				}
			}
			if i != fieldNums-1 {
				query += ", "
			}
			log.Println(query)
		}
	}
	return createTableQueryPrefix + query + querySuffix, nil
}

// Creates a query to create a table based on the struct that is passed in.
// Returns an error if the struct is empty/invalid.
func createTableCreationQueryV2(object reflect.Value, tableName string, pkFieldName string, notNullFields []string) (string, error) {
	createTableQueryPrefix := "CREATE TABLE " + tableName + " (\n"
	querySuffix := "\n);"
	queryParts := []string{}
	fieldNums := object.NumField()

	if fieldNums == 0 {
		return createTableQueryPrefix + querySuffix, nil
	}

	notNullMap := make(map[string]bool)
	for _, field := range notNullFields {
		notNullMap[field] = true
	}

	for i := 0; i < fieldNums; i++ {
		field := object.Type().Field(i)
		fieldType := field.Type
		fieldName := field.Name

		postgresType, ok := goTypeToPostgresType[fieldType.Name()]
		if !ok {
			return "", errors.New("unsupported field type: " + fieldType.Name())
		}

		columnDef := "\t" + fieldName

		if fieldName == pkFieldName {
			columnDef += " SERIAL PRIMARY KEY"
		} else {
			columnDef += " " + postgresType
			switch postgresType {
			case "VARCHAR":
				columnDef += "(255)"
			case "TEXT":
				columnDef += "(65535)"
			case "DATE":
				columnDef += "(10)"
			case "TIME":
				columnDef += "(8)"
			}
			if notNullMap[fieldName] {
				columnDef += " NOT NULL"
			}
		}
		queryParts = append(queryParts, columnDef)
	}

	query := createTableQueryPrefix + join(queryParts, ",\n") + querySuffix
	return query, nil
}

// Returns true if only one PK is defined in the struct,
// false if more than one Pk is defined or no PK is defined
func determinePrimaryKey(object reflect.Value) (string, error) {
	var pkOccurences = 0
	var fieldName = ""
	fieldNums := object.NumField()
	for i := 0; i < fieldNums; i++ {
		field := object.Type().Field(i)
		if tableName := field.Tag.Get("dapFieldAttrs"); tableName != "" {
			if strings.Contains(strings.ToLower(tableName), "pk") {
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
	return "", fmt.Errorf("error while determining PK")
}

// Returns a list of field names that are marked as NOT NULL in the struct
func determineNotNullFields(object reflect.Value) ([]string, error) {
	var notNullFields []string
	fieldNums := object.NumField()
	for i := 0; i < fieldNums; i++ {
		field := object.Type().Field(i)
		if tableName := field.Tag.Get("dapFieldAttrs"); tableName != "" {
			log.Printf("Field: %v, Attrs: %v", field.Name, tableName)
			if strings.Contains(tableName, "NOT NULL") {
				notNullFields = append(notNullFields, field.Name)
			}
		}
	}
	return notNullFields, nil
}

// Joins the elements of a slice of strings with a separator.
func join(elems []string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0]
	default:
		n := len(sep) * (len(elems) - 1)
		for i := 0; i < len(elems); i++ {
			n += len(elems[i])
		}

		b := make([]byte, n)
		bp := copy(b, elems[0])
		for _, s := range elems[1:] {
			bp += copy(b[bp:], sep)
			bp += copy(b[bp:], s)
		}
		return string(b)
	}
}
