package main

import (
	"busline/platform/db"

	"database/sql"
	"log"
	"os"
)

func main() {

	dbConn, err := setupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	err = db.RunMigrations(dbConn)

	if err != nil {
		log.Fatal(err)
	}
}

/* initializes the connection to postgres
 */
func setupDatabase() (*sql.DB, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}
	return db, nil
}
