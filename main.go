package main

import (
	"busline/cmd/cron"
	"busline/cmd/server"
	"busline/internal/app"
	"busline/internal/gtfs"
	"busline/platform/db"

	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	dbConn, err := setupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	err = db.RunMigrations(dbConn)

	if err != nil {
		log.Fatal(err)
	}

	gtfsRepository := gtfs.NewRepository(dbConn)

	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(gtfsRepository)

	cron.StartCronJobs(app)

	server.StartServer(ctx, app)

}

/* initializes the connection to postgres
 */
func setupDatabase() (*sql.DB, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}
	return db, nil
}
