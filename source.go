package belt

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mitchellh/mapstructure"
)

type dbRow map[string]string

const (
	nullRaw = "\\N"
)

// dbSource format: "USER/PASSWORD@tcp(HOST:PORT)/DATABASE"
func getFromMySql(dbSource string, output interface{}, query string, args ...interface{}) []I {
	db, err := sql.Open("mysql", dbSource)
	errorHandlerWithFatal(err)
	defer db.Close()
	rows, err := db.Query(query, args...)
	errorHandlerWithFatal(err)
	defer rows.Close()

	columns, err := rows.Columns()
	errorHandlerWithFatal(err)

	rawResult := make([][]byte, len(columns))
	results := make([]I, 0)
	dest := make([]interface{}, len(columns))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		err := rows.Scan(dest...)
		errorHandlerWithFatal(err)
		row := make(dbRow)
		for i, raw := range rawResult {
			if raw == nil {
				row[columns[i]] = nullRaw
			} else {
				row[columns[i]] = string(raw)
			}
		}
		mapstructure.Decode(row, &output)
		results = append(results, output)
	}
	return results
}
