package dapter

var goTypeToPostgresType = map[string]string{
	"int":       "INTEGER",
	"int8":      "INTEGER",
	"int16":     "INTEGER",
	"int32":     "INTEGER",
	"int64":     "BIGINT",
	"uint":      "INTEGER",
	"uint8":     "INTEGER",
	"uint16":    "INTEGER",
	"uint32":    "INTEGER",
	"uint64":    "BIGINT",
	"float32":   "REAL",
	"float64":   "DOUBLE PRECISION",
	"bool":      "BOOLEAN",
	"string":    "VARCHAR",
	"time.Time": "TIMESTAMP",
	"[]byte":    "BYTEA",
}