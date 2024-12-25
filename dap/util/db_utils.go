package util

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/piyushpatil22/dapter/log"
)

func GetDBMaxRetryCounts() int {
	reconnectCount := 5
	dbReconnectMaxCount := os.Getenv("DB_RECONNECT_MAX_COUNT")
	if dbReconnectMaxCount != "" {
		parsedCount, err := strconv.ParseInt(dbReconnectMaxCount, 10, 64)
		if err != nil {
			log.Log.Err(err).Msg("Error parsing db reconnect max count")
		} else {
			reconnectCount = int(parsedCount)
		}
	}
	return reconnectCount
}

func ConnectToDatabase() (*sql.DB, error) {
	log.Log.Info().Msg("Connecting to db...")
	//db conn
	dbConn, err := NewPostgreSQLStore(DBConfig{
		DBAddress:  "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "root",
		DBName:     "traders",
	})
	if err != nil {
		log.Log.Err(err).Msg("Error connecting to db")
		return nil, err
	}
	log.Log.Info().Msg("Pinging to db...")
	err = dbConn.Ping()
	if err != nil {
		log.Log.Err(err).Msg("Error pinging to db")
		return nil, err
	}
	return dbConn, nil
}

func NewPostgreSQLStore(cfg DBConfig) (*sql.DB, error) {
	connection := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", cfg.DBAddress, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB, tableName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, created_at TIMESTAMP, updated_at TIMESTAMP)", tableName))
	if err != nil {
		return err
	}
	return nil
}



func GetFields(ent interface{}, nested bool) []string {
	fields := make([]string, 0)
	for i := 0; i < reflect.TypeOf(ent).NumField(); i++ {
		if nested {
			if reflect.TypeOf(ent).Field(i).Name == "Base" {
				for j := 0; j < reflect.ValueOf(ent).Field(i).NumField(); j++ {
					fields = append(fields, reflect.TypeOf(ent).Field(i).Type.Field(j).Tag.Get("json"))
				}
			}
		}
		if reflect.TypeOf(ent).Field(i).Name == "Base" {
			continue
		}
		fields = append(fields, reflect.TypeOf(ent).Field(i).Tag.Get("json"))
	}
	return fields
}

func GetFieldsWithTypes(ent interface{}) map[string]reflect.Kind {
	fields := make(map[string]reflect.Kind)
	for i := 0; i < reflect.TypeOf(ent).NumField(); i++ {
		fields[reflect.TypeOf(ent).Field(i).Tag.Get("json")] = reflect.TypeOf(ent).Field(i).Type.Kind()
	}
	return fields
}

func SanitizeQuery(query string) (string, error) {
	if !strings.Contains(query, "INSERT INTO") || !strings.Contains(query, "VALUES") {
		return "", errors.New("invalid query: must contain 'INSERT INTO' and 'VALUES'")
	}

	parts := strings.Split(query, "VALUES")
	if len(parts) != 2 {
		return "", errors.New("invalid query format")
	}

	columnsPart := strings.TrimSpace(parts[0])
	valuesPart := strings.TrimSpace(parts[1])

	columnsStart := strings.Index(columnsPart, "(")
	columnsEnd := strings.LastIndex(columnsPart, ")")
	valuesStart := strings.Index(valuesPart, "(")
	valuesEnd := strings.LastIndex(valuesPart, ")")

	if columnsStart == -1 || columnsEnd == -1 || valuesStart == -1 || valuesEnd == -1 {
		return "", errors.New("invalid query format: missing parentheses")
	}

	columns := strings.Split(columnsPart[columnsStart+1:columnsEnd], ",")
	values := strings.Split(valuesPart[valuesStart+1:valuesEnd], ",")

	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}

	if columns[len(columns)-1] == "" {
		columns = columns[:len(columns)-1]
	}
	if values[len(values)-1] == "" {
		values = values[:len(values)-1]
	}

	if len(columns) != len(values) {
		return "", errors.New("mismatch between the number of columns and values")
	}

	cleanedColumns := strings.Join(columns, ", ")
	cleanedValues := strings.Join(values, ", ")

	cleanedQuery := fmt.Sprintf(
		"%s (%s) VALUES (%s);",
		strings.TrimSpace(columnsPart[:columnsStart]),
		cleanedColumns,
		cleanedValues,
	)

	return cleanedQuery, nil
}
