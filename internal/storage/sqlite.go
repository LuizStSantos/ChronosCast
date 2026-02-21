package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS streams (
		id TEXT PRIMARY KEY, entrada TEXT, saida TEXT, dias TEXT, inicio TEXT, fim TEXT
	);`
	_, err = db.Exec(sqlStmt)
	return db, err
}
