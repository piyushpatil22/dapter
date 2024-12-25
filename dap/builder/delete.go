package builder

import "fmt"

func generateDeleteQuery(builder *QueryBuilder) (string, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = %s", builder.TableName, builder.Value.Field(0).Interface().(string))
	return query, nil
}
