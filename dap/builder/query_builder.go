package builder

import (
	"fmt"
	"reflect"

	"github.com/piyushpatil22/dapter/dap/filter"
	"github.com/piyushpatil22/dapter/log"
)

const (
	INSERT = "INSERT"
	UPDATE = "UPDATE"
	DELETE = "DELETE"
	GET    = "GET"
)

var ENTITY_TABLE_NAME_MAPPING = map[string]string{
	"Instrument": "instruments",
	"User":       "users",
}

type QueryBuilder struct {
	Output    interface{}
	Value     reflect.Value
	Type      reflect.Type
	TableName string
	Filters   []filter.Filter
	argCount  int
	tableName string
	query     string
	args      []interface{}
}

func (qb *QueryBuilder) AppendFilters(filters []filter.Filter) {
	qb.Filters = append(qb.Filters, filters...)
}

func NewQueryBuilder(ent interface{}) *QueryBuilder {
	return &QueryBuilder{
		Value:     reflect.ValueOf(ent),
		Type:      reflect.TypeOf(ent),
		TableName: GetTableName(ent),
		Filters:   make([]filter.Filter, 0),
	}
}

func GenerateQuery(ent interface{}, tableName string, queryType string, filters []filter.Filter) (string, error) {
	builder := NewQueryBuilder(ent)
	switch queryType {
	case INSERT:
		return generateInsertQuery(builder)
	case UPDATE:
		return generateUpdateQuery(builder)
	case DELETE:
		return generateDeleteQuery(builder)
	case GET:
		if filters != nil {
			builder.AppendFilters(filters)
		}
		return generateGetQuery(builder)
	default:
		return "", fmt.Errorf("invalid query type %s", queryType)
	}
}

func isZeroValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Interface() == ""
	case reflect.Bool:
		return value.Interface() == false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Interface() == int64(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Interface() == uint64(0)
	case reflect.Float32, reflect.Float64:
		return value.Interface() == float64(0)
	case reflect.Complex64, reflect.Complex128:
		return value.Interface() == complex128(0)
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr:
		return value.IsNil()
	case reflect.Struct:
		return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
	default:
		return false
	}
}

func GetTableName(ent interface{}) string {
	entType := reflect.TypeOf(ent)
	entName := entType.Name()
	if val, ok := ENTITY_TABLE_NAME_MAPPING[entName]; ok {
		return val
	} else {
		return ""
	}
}
func GetTableName2(ent interface{}) string {
	entName := ""
	//if slice, get the type of the slice
	if reflect.TypeOf(ent).Kind() == reflect.Slice {
		entName = reflect.TypeOf(ent).Elem().Name()
	} else {
		entName = reflect.TypeOf(ent).Name()
	}
	log.Log.Info().Str("entity", entName).Msg("Entity")
	if val, ok := ENTITY_TABLE_NAME_MAPPING[entName]; ok {
		return val
	} else {
		return ""
	}
}
