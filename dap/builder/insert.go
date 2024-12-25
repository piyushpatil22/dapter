package builder

import (
	"fmt"
	"time"

	"github.com/piyushpatil22/dapter/dap/util"
)

func generateInsertQuery(builder *QueryBuilder) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s (", builder.TableName)
	for i := 0; i < builder.Value.NumField(); i++ {
		//check if the field is Base struct
		if builder.Type.Field(i).Name == "Base" {
			for j := 0; j < builder.Value.Field(i).NumField(); j++ {
				if builder.Value.Field(i).Type().Field(j).Name == "ID" {
					continue
				}

				query += builder.Value.Field(i).Type().Field(j).Tag.Get("json")
				if j != builder.Value.Field(i).NumField() {
					query += ", "
				}
			}
			continue
		}
		if isZeroValue(builder.Value.Field(i)) || builder.Type.Field(i).Tag.Get("json") == "id" {
			continue
		}
		query += builder.Type.Field(i).Tag.Get("json")
		if i != builder.Value.NumField()-1 {
			query += ", "
		}
	}
	query += ") VALUES ("
	for i := 0; i < builder.Value.NumField(); i++ {
		if builder.Type.Field(i).Name == "Base" {
			for j := 0; j < builder.Value.Field(i).NumField(); j++ {
				if builder.Value.Field(i).Type().Field(j).Name == "ID" {
					continue
				}
				if builder.Value.Field(i).Type().Field(j).Name == "CreatedAt" || builder.Value.Field(i).Type().Field(j).Name == "UpdatedAt" {
					currentTime := time.Now().Format("2006-01-02 15:04:05")
					query += fmt.Sprintf("'%v'", currentTime)
				} else {
					query += fmt.Sprintf("'%v'", builder.Value.Field(i).Field(j).Interface())
				}
				if j != builder.Value.Field(i).NumField() {
					query += ", "
				}
			}
			continue
		}
		if isZeroValue(builder.Value.Field(i)) || builder.Type.Field(i).Tag.Get("json") == "id" {
			continue
		}
		query += fmt.Sprintf("'%v'", builder.Value.Field(i).Interface())
		if i != builder.Value.NumField()-1 {
			query += ", "
		}
	}
	query += ")"
	cleanedQuery, err := util.SanitizeQuery(query)
	if err != nil {
		return "", err
	}
	return cleanedQuery, nil
}