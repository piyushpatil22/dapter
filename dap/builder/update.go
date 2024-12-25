package builder

import (
	"fmt"
	"time"
)

func generateUpdateQuery(builder *QueryBuilder) (string, error) {
	query := fmt.Sprintf("UPDATE %s SET ", builder.TableName)
	var id string
	for i := 0; i < builder.Value.NumField(); i++ {
		//check if the field is Base struct
		if builder.Type.Field(i).Name == "Base" {
			for j := 0; j < builder.Value.Field(i).NumField(); j++ {
				if builder.Value.Field(i).Type().Field(j).Name == "ID" {
					id = builder.Value.Field(i).Field(j).Interface().(string)
					if id == "" {
						return "", fmt.Errorf("id is empty")
					}
					continue
				}
				if isZeroValue(builder.Value.Field(i).Field(j)) {
					continue
				}
				if builder.Value.Field(i).Type().Field(j).Name == "UpdatedAt" {
					currentTime := time.Now().Format("2006-01-02 15:04:05")
					query += fmt.Sprintf("%s = '%v'", builder.Value.Field(i).Type().Field(j).Tag.Get("json"), currentTime)
				} else {
					query += fmt.Sprintf("%s = '%v'", builder.Value.Field(i).Type().Field(j).Tag.Get("json"), builder.Value.Field(i).Field(j).Interface())
				}
				if j != builder.Value.Field(i).NumField() {
					query += ", "
				}
			}
			continue
		}
		if isZeroValue(builder.Value.Field(i)) {
			continue
		}
		query += fmt.Sprintf("%s = '%v'", builder.Type.Field(i).Tag.Get("json"), builder.Value.Field(i).Interface())
		if i != builder.Value.NumField()-1 {
			query += ", "
		}
	}
	query += fmt.Sprintf(" WHERE id = %s", id)
	return query, nil
}
