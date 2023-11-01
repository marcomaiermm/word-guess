package database

import (
	"database/sql"
	"fmt"
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
	Seed()

	return nil
}

func Seed() error {
	row := DB.QueryRow("SELECT 1 FROM word LIMIT 1")

	var exists int
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// if there are no words in the database, seed it.
	// Get the words from the file data.txt. The first line is a comment in the file, so we skip it.
	// Split the words by new line and insert them into the database.
	const BATCH_SIZE = 500
	if err == sql.ErrNoRows {
		words, err := os.ReadFile("pkg/database/data.txt")
		if err != nil {
			return err
		}
		lines := strings.Split(string(words), "\n")[1:] // Skip the first line

		for i := 0; i < len(lines); i += BATCH_SIZE {
			endIndex := i + BATCH_SIZE
			if endIndex > len(lines) {
				endIndex = len(lines)
			}

			valueStrings := []string{}
			valueArgs := []interface{}{}
			for _, line := range lines[i:endIndex] {
				valueStrings = append(valueStrings, "(?)")
				valueArgs = append(valueArgs, line)
			}
			stmt := fmt.Sprintf("INSERT INTO word (word) VALUES %s", strings.Join(valueStrings, ","))

			_, err = DB.Exec(stmt, valueArgs...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
