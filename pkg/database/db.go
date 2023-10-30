package database

import (
	"database/sql"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(url string) error {
	db, err := sql.Open("sqlite", url)
	if err != nil {
		return err
	}

	scripts, err := os.ReadFile("pkg/database/scripts/tables.sql")
	if err != nil {
		return err
	}

	statements := strings.Split(string(scripts), ";")

	for _, statement := range statements {
		if strings.TrimSpace(statement) != "" {
			_, err := db.Exec(statement)
			if err != nil {
				return err
			}
		}
	}

	DB = db

	return nil
}
