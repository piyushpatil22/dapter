package dap

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/piyushpatil22/dapter/dap/builder"
	"github.com/piyushpatil22/dapter/dap/executor"
	"github.com/piyushpatil22/dapter/dap/filter"
	"github.com/piyushpatil22/dapter/dap/parser"
	"github.com/piyushpatil22/dapter/dap/util"
	"github.com/piyushpatil22/dapter/log"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

func (s *Store) Close() error {
	return s.DB.Close()
}

func (s *Store) Insert(ent any) error {
	ok, err := s.checkTableExists(builder.GetTableName(ent))
	if err != nil {
		log.Log.Err(err).Msg("Error checking table exists")
		return err
	}
	if !ok {
		err = s.CreateTable(ent)
		if err != nil {
			log.Log.Err(err).Msg("Error creating table")
			return err
		}
	} else {
		new_cols, err := s.CheckTableNeedsUpdate(ent)
		if err != nil {
			log.Log.Err(err).Msg("Error checking table needs update")
			return err
		}
		log.Log.Info().Interface("new_cols", new_cols).Msg("New columns")
		if len(new_cols) > 0 {
			log.Log.Info().Msg("Table needs update")
			err = s.UpdateTable(ent)
			if err != nil {
				log.Log.Err(err).Msg("Error updating table")
				return err
			}
		}
	}
	//get the type of the entity
	entType := reflect.TypeOf(ent)
	entName := entType.Name()
	db_table_name := builder.GetTableName(ent)
	if db_table_name == "" {
		log.Log.Err(err).Msg("Table name not found for entity")
		return fmt.Errorf("table name not found for entity %s", entName)
	}
	//get the fields of the entity
	//create a query to insert the entity
	query, err := builder.GenerateQuery(ent, db_table_name, builder.INSERT, nil)
	if err != nil {
		log.Log.Err(err).Msg("Error generating insert query")
	} else {
		log.Log.Info().Str("query", query).Msg("Query")
	}

	//execute the query
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Log.Err(err).Msg("Error executing query")
		return err
	}
	defer rows.Close()
	log.Log.Info().Msg("Query executed")
	return nil
}

func (s *Store) Update(ent any) error {
	//get the type of the entity
	entType := reflect.TypeOf(ent)
	entName := entType.Name()
	db_table_name := builder.GetTableName(ent)
	if db_table_name == "" {
		log.Log.Error().Msg("Table name not found for entity")
		return fmt.Errorf("table name not found for entity %s", entName)
	}
	//get the fields of the entity
	//create a query to update the entity
	query, err := builder.GenerateQuery(ent, db_table_name, builder.UPDATE, nil)
	if err != nil {
		log.Log.Err(err).Msg("Error generating update query")
	} else {
		log.Log.Info().Str("query", query).Msg("Query")
	}

	//execute the query
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Log.Err(err).Msg("Error executing query")
		return nil
	}
	defer rows.Close()
	log.Log.Info().Msg("Query executed")

	return nil
}

func (s *Store) Delete(ent any) error {
	//get the type of the entity
	entType := reflect.TypeOf(ent)
	entName := entType.Name()
	db_table_name := builder.GetTableName(ent)
	if db_table_name == "" {
		log.Log.Error().Msg("Table name not found for entity")
		return fmt.Errorf("table name not found for entity %s", entName)
	}
	//get the fields of the entity
	//create a query to delete the entity
	query, err := builder.GenerateQuery(ent, db_table_name, builder.DELETE, nil)
	if err != nil {
		log.Log.Err(err).Msg("Error generating delete query")
	} else {
		log.Log.Info().Str("query", query).Msg("Query")
	}

	//execute the query

	return nil
}

func (s *Store) BulkInsert(ent []interface{}) error {
	for _, e := range ent {
		err := s.Insert(e)
		if err != nil {
			log.Log.Err(err).Msg("Error inserting entity")
			return err
		}
	}
	return nil
}

func (s *Store) GetByFilter(result interface{}, filters filter.Filter, ent interface{}) error {
	if ent == nil {
		log.Log.Error().Msg("Entity is nil")
		return fmt.Errorf("entity is nil")
	}
	if !validateModelsMatch(result, ent) {
		log.Log.Error().Msg("Error validating models match")
		return fmt.Errorf("models do not match")
	}
	entType := reflect.TypeOf(ent)
	db_table_name := builder.GetTableName(ent)
	entName := entType.Name()
	if db_table_name == "" {
		log.Log.Error().Msg("Table name not found for entity")
		return fmt.Errorf("table name not found for entity %s", entName)
	}
	//get the fields of the entity
	//create a query to get the entity by id
	query, err := builder.GenerateQuery(ent, db_table_name, builder.GET, []filter.Filter{filters})
	if err != nil {
		log.Log.Err(err).Msg("Error generating get query")
	} else {
		log.Log.Info().Str("query", query).Msg("Query")
	}

	//execute the query
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Log.Err(err).Msg("Error executing query")
		return err
	}
	defer rows.Close()
	log.Log.Info().Msg("Query executed")
	rowsData := ConvertToDapRow2(rows)
	if len(rowsData) == 0 {
		return ErrNoRowsFound
	}
	err = parser.Parse2Struct(result, rowsData)
	if err != nil {
		log.Log.Err(err).Msg("Error parsing rows to struct")
		return err
	}

	return nil
}

// func (s *Store) GetByID(ent any, id interface{}) ([]parser.DapRow, error) {
// 	filter := filter.Filter{
// 		Field: "id",
// 		Value: id,
// 	}
// 	return s.GetByFilter(ent, filter, nil)
// }

func (s *Store) CreateTable(ent any) error {
	query, err := executor.CreateEntityTableWithFields(ent)
	if err != nil {
		log.Log.Err(err).Msg("Error creating table")
		return err
	}
	log.Log.Info().Str("query", query).Msg("Query")

	//execute the query
	result, err := s.DB.Exec(query)
	if err != nil {
		log.Log.Err(err).Msg("Error executing query")
		return err
	}
	log.Log.Info().Interface("result", result).Msg("Table created")

	return nil
}

func (s *Store) FetchColumns(tableName string) ([]string, error) {
	rows, err := s.DB.Query(fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s' ORDER BY ordinal_position", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns := make([]string, 0)
	for rows.Next() {
		var column string
		err := rows.Scan(&column)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, nil
}

func (s *Store) CheckTableNeedsUpdate(ent any) ([]string, error) {
	db_table_name := builder.GetTableName(ent)
	log.Log.Info().Str("table_name", db_table_name).Msg("Table name")
	DBcols, err := s.FetchColumns(db_table_name)
	if err != nil {
		return nil, err
	}
	objCols := util.GetFields(ent, true)

	log.Log.Info().Interface("DBcols", DBcols).Msg("DB columns")
	log.Log.Info().Interface("objCols", objCols).Msg("Object columns")
	var newCols []string
	for _, col := range objCols {
		found := false
		for _, dbcol := range DBcols {
			if col == dbcol {
				found = true
				break
			}
		}
		if !found {
			log.Log.Info().Str("col", col).Msg("Column not found in DB")
			newCols = append(newCols, col)
		}
	}
	return newCols, nil
}

func (s *Store) UpdateTable(ent any) error {
	newCols, err := s.CheckTableNeedsUpdate(ent)
	if err != nil {
		log.Log.Err(err).Msg("Error checking table needs update")
		return err
	}
	if len(newCols) == 0 {
		log.Log.Info().Msg("Table does not need update")
		return nil
	}
	fieldsWithType := util.GetFieldsWithTypes(ent)
	query := fmt.Sprintf("ALTER TABLE %s ", builder.GetTableName(ent))
	for i, col := range newCols {
		//determine the type of the column
		typeOfCol := fieldsWithType[col]
		var DBType string
		switch typeOfCol {
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
			log.Log.Err(err).Msg("Unsupported field type")
			return fmt.Errorf("unsupported field type : %+v", typeOfCol)
		}
		query += fmt.Sprintf("ADD COLUMN %s %s", col, DBType)
		if i != len(newCols)-1 {
			query += ", "
		}

	}
	log.Log.Info().Str("query", query).Msg("Query")
	result, err := s.DB.Exec(query)
	if err != nil {
		log.Log.Err(err).Msg("Error executing query")
		return err
	}
	log.Log.Info().Interface("result", result).Msg("Table updated")
	return nil
}

func (s *Store) checkTableExists(tableName string) (bool, error) {
	rows, err := s.DB.Query(fmt.Sprintf("SELECT to_regclass('public.%s')", tableName))
	if err != nil {
		return false, err
	}
	defer rows.Close()
	var exists bool
	var row interface{}
	for rows.Next() {
		err := rows.Scan(&row)
		if err != nil {
			return false, err
		}
	}
	if row == nil {
		exists = false
	} else {
		exists = true
	}

	return exists, nil
}

func ConvertToDapRow2(rows *sql.Rows) []parser.DapRow {
	outputRows := make([]parser.DapRow, 0)
	columns, err := rows.Columns()
	if err != nil {
		log.Log.Err(err).Msg("Error getting columns")
		return outputRows
	}
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Log.Err(err).Msg("Error scanning row")
			continue
		}
		row := parser.DapRow{
			Columns: columns,
			Values:  make([]interface{}, len(columns)),
		}
		copy(row.Values, values)
		outputRows = append(outputRows, row)
	}
	return outputRows
}

func validateModelsMatch(ent1, ent2 any) bool {
	ent1Type := reflect.TypeOf(ent1)
	if ent1Type.Kind() != reflect.Ptr || ent1Type.Elem().Kind() != reflect.Slice {
		return false
	}
	sliceElemType := ent1Type.Elem().Elem()
	ent2Type := reflect.TypeOf(ent2)
	return sliceElemType == ent2Type
}
