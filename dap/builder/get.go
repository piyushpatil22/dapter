package builder

import "fmt"

func generateGetQuery(builder *QueryBuilder) (string, error) {
	if builder.Filters != nil && len(builder.Filters) > 0 {
		query := fmt.Sprintf("SELECT * FROM %s WHERE ", builder.TableName)
		for i, filter := range builder.Filters {
			query += fmt.Sprintf("%s = '%s'", filter.Field, filter.Value)
			if i != len(builder.Filters)-1 {
				query += " AND "
			}
		}
		return query, nil
	} else {
		return "", fmt.Errorf("no filters provided")
	}
}
