# Platform

This document outlines how to interact with platforms that the application relies on

## DB

The app uses postgres as its database.

### Migrations Setup

golang-migrate is the migration tool the application uses and runs on app start.

download it by running:
`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### Updating Tables

In the current directory, run to create a setup of sql migration files:
`migrate create -ext sql -dir ./platform/db/migrations -seq "name of migration"`
