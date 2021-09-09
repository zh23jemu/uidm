package db

import (
	"database/sql"
	"uidm/output"

	_ "github.com/denisenkom/go-mssqldb"
)

func QueryMsSql(sqlquery string) *sql.Rows {
	conn, err := sql.Open("mssql", connString)
	output.GenerateLog(err, "DB connection", connString, false)
	defer conn.Close()

	rows, err := conn.Query(sqlquery)
	output.GenerateLog(err, "DB query", sqlquery,  false)

	return rows
}

func QueryMsSqlRow(sqlquery string) (count int) {
	conn, err := sql.Open("mssql", connString)
	output.GenerateLog(err, "DB connection", connString,  false)
	defer conn.Close()

	err = conn.QueryRow(sqlquery).Scan(&count)
	output.GenerateLog(err, "DB query", sqlquery,  false)

	return count
}
