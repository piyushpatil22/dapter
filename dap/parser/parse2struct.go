package parser

import (
	"errors"
	"fmt"
	"reflect"
)

// TODO need to fix error returns, define error in 1 place
type DapRow struct {
	Columns []string
	Values  []interface{}
}

func Parse2Struct(result interface{}, rows []DapRow) error {
	resultVal := reflect.ValueOf(result)
	if resultVal.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer")
	}
	elemType := resultVal.Elem().Type()
	if elemType.Kind() == reflect.Slice {
		sliceElemType := elemType.Elem()
		sliceVal := reflect.MakeSlice(elemType, len(rows), len(rows))
		for i := 0; i < len(rows); i++ {
			elem := reflect.New(sliceElemType).Elem()
			if err := fillStruct(elem, rows[i]); err != nil {
				return err
			}
			sliceVal.Index(i).Set(elem)
		}

		resultVal.Elem().Set(sliceVal)
		return nil
	}
	if len(rows) > 0 {
		return fillStruct(resultVal.Elem(), rows[0])
	}
	return errors.New("no rows to parse")
}

func fillStruct(structVal reflect.Value, row DapRow) error {
	if structVal.Kind() != reflect.Struct {
		return errors.New("target value must be a struct")
	}
	for i, column := range row.Columns {
		field := structVal.FieldByNameFunc(func(name string) bool {
			field, _ := structVal.Type().FieldByName(name)
			tag := field.Tag.Get("json")
			return tag == column || name == column
		})
		if !field.IsValid() {
			for j := 0; j < structVal.NumField(); j++ {
				embeddedField := structVal.Field(j)
				if embeddedField.Kind() == reflect.Struct {
					embeddedStruct := structVal.Type().Field(j).Name
					field = structVal.FieldByName(embeddedStruct).FieldByName(column)
					if field.IsValid() && field.CanSet() {
						break
					}
				}
			}
		}
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		value := row.Values[i]
		if value == nil {
			field.Set(reflect.Zero(field.Type()))
			continue
		}
		val := reflect.ValueOf(value)
		if val.Type().AssignableTo(field.Type()) {
			field.Set(val)
		} else if val.Type().ConvertibleTo(field.Type()) {
			field.Set(val.Convert(field.Type()))
		} else {
			return fmt.Errorf("cannot assign value of type %s to field %s of type %s", val.Type(), column, field.Type())
		}
	}
	return nil
}
