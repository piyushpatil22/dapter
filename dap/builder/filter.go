package builder

import (
	"fmt"
	"strings"
)

//TODO make constants for SQL keywords and logical operators

type Condition interface {
	IsCondition()
}
type LogicalCondition struct {
	Operator   string
	Conditions []Condition
}
type QueryCondition struct {
	Field    string
	Value    string
	Operator string
}

func (qc QueryCondition) IsCondition()   {}
func (qc LogicalCondition) IsCondition() {}

func Select(result interface{}, ent interface{}) *QueryBuilder {
	return &QueryBuilder{
		tableName: GetTableName(ent),
		query:     "SELECT * FROM " + GetTableName2(ent),
		Output:    result,
	}
}

func (qb *QueryBuilder) WHERE(conditions []Condition) *QueryBuilder {
	qb.query += " WHERE "
	qb.query += qb.buildConditions(conditions)
	return qb
}
func (qb *QueryBuilder) GetQueryArgs() ([]interface{}, string) {
	return qb.args, qb.query
}
func (qb *QueryBuilder) buildConditions(conditions []Condition) string {
	var parts []string
	var args []interface{}
	for _, con := range conditions {
		switch conType := con.(type) {
		case QueryCondition:
			qb.argCount += 1
			parts = append(parts, fmt.Sprintf("%s %s $%v", conType.Field, conType.Operator, qb.argCount))
			args = append(args, conType.Value)
		case LogicalCondition:
			subQuery := qb.buildConditions(conType.Conditions)
			parts = append(parts, fmt.Sprintf("(%s)", subQuery))
		}
	}
	qb.args = append(qb.args, args...)
	return strings.Join(parts, " AND ")
}
