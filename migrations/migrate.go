package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

// Migrate runs the database migrations
func Migrate(db *sql.DB) error {
	queries := []struct {
		fileName string
		fn       func(*sql.DB, string) error
	}{
		{"queries/create_all_tables.sql", runSQLFile},
		// {"queries/create_users_table.sql", runSQLFile},
		// {"queries/create_otps_table.sql", runSQLFile},
	}

	for _, q := range queries {
		if err := q.fn(db, q.fileName); err != nil {
			return err
		}
	}

	fmt.Println("👍 Migration complete")
	return nil
}

func runSQLFile(db *sql.DB, fileName string) error {
	query, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error reading SQL file %s: %v", fileName, err)
		return err
	}

	_, err = db.Exec(string(query))
	if err != nil {
		log.Printf("Error executing SQL file %s: %v", fileName, err)
		return err
	}
	return nil
}
