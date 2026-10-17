package main

import (
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

	madisonGtfsConfig := gtfs.GTFSConfig{StaticUrl: "https://transitdata.cityofmadison.com/GTFS/mmt_gtfs.zip", City: gtfs.Madison}
	madisonGtfs, err := gtfs.NewGTFS(madisonGtfsConfig)
	gtfsRepository := gtfs.NewRepository(dbConn)

	if err != nil {
		log.Fatal(err)
	}

	err = madisonGtfs.DownloadStatic(ctx)
	if err != nil {
		log.Fatal(err)
	}

	madisonGtfs.ParseStatic()

	feedExists, err := gtfsRepository.DoesFeedExist(madisonGtfs.StaticHash)

	if err != nil {
		log.Fatal("gtfs: %w", err)
	}

	if !feedExists {
		err = gtfsRepository.SaveStaticGTFS(ctx, madisonGtfs)
		if err != nil {
			log.Fatal(err)
		}

		err = madisonGtfs.StreamAndSaveStopTimes(ctx, gtfsRepository)
		if err != nil {
			log.Fatal(err)
		}
	}
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
