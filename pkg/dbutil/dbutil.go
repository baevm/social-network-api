package dbutil

import (
	"fmt"
	"strconv"
	"strings"
)

// https://github.com/jackc/pgx/issues/764#issuecomment-685249471
func GetBulkInsertSQL(SQLString string, rowValueSQL string, numRows int) string {
	// Combine the base SQL string and N value strings
	valueStrings := make([]string, 0, numRows)
	for i := 0; i < numRows; i++ {
		valueStrings = append(valueStrings, "("+rowValueSQL+")")
	}
	allValuesString := strings.Join(valueStrings, ",")
	SQLString = fmt.Sprintf(SQLString, allValuesString)

	// Convert all of the "?" to "$1", "$2", "$3", etc.
	// (which is the way that pgx expects query variables to be)
	numArgs := strings.Count(SQLString, "?")
	SQLString = strings.ReplaceAll(SQLString, "?", "$%v")
	numbers := make([]interface{}, 0, numRows)
	for i := 1; i <= numArgs; i++ {
		numbers = append(numbers, strconv.Itoa(i))
	}
	return fmt.Sprintf(SQLString, numbers...)
}

// GetBulkInsertSQLString is a helper function to prepare a SQL query for a bulk insert.
// GetBulkInsertSQLString is used over getBulkInsertSQL when all of the values are plain question
// marks (e.g. a 1-for-1 value insertion).
// The example given for getBulkInsertSQL is such a query.
func GetBulkInsertSQLString(SQLString string, numArgsPerRow int, numRows int) string {
	questionMarks := make([]string, 0, numArgsPerRow)
	for i := 0; i < numArgsPerRow; i++ {
		questionMarks = append(questionMarks, "?")
	}
	rowValueSQL := strings.Join(questionMarks, ", ")
	return GetBulkInsertSQL(SQLString, rowValueSQL, numRows)
}
