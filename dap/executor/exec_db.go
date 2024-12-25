package executor

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/piyushpatil22/dapter/dap/builder"
)

func CreateEntityTable(db *sql.DB, ent interface{}) error {
	tableName := builder.GetTableName(ent)
	if tableName == "" {
		return fmt.Errorf("table name not found for entity")
	}
	_, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, created_at TIMESTAMP, updated_at TIMESTAMP)", tableName))
	if err != nil {
		return err
	}
	return nil
}

func CreateEntityTableWithFields(ent interface{}) (string, error) {
	tableName := builder.GetTableName(ent)
	if tableName == "" {
		return "", fmt.Errorf("table name not found for entity")
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, created_at TIMESTAMP, updated_at TIMESTAMP", tableName)
	for i := 0; i < reflect.TypeOf(ent).NumField(); i++ {
		//check field type like string, int, float, etc
		var DBType string
		switch reflect.TypeOf(ent).Field(i).Type.Kind() {
		case reflect.String:
			DBType = "TEXT"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			DBType = "INT"
		case reflect.Float32, reflect.Float64:
			DBType = "FLOAT"
		case reflect.Bool:
			DBType = "BOOLEAN"
		case reflect.Struct:
			continue
		default:
			return "", fmt.Errorf("unsupported field type : %+v", reflect.TypeOf(ent).Field(i).Type.Kind())
		}
		query += fmt.Sprintf(", %s %s", reflect.TypeOf(ent).Field(i).Tag.Get("json"), DBType)
	}
	query += ")"
	return query, nil
}
