package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/defryheryanto/ai-assistant/config"
	_ "github.com/golang-migrate/migrate/source"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
)

func main() {
	downFlag := flag.Bool("down", false, "database migration down")
	flag.Parse()

	config.Init()

	log.Println("Opening database connection")

	db, err := sql.Open("postgres", config.DatabaseConnectionString)
	log.Println("Database connected.")
	if err != nil {
		log.Fatalf("error opening migration database - %v\n", err)
		return
	}

	log.Println("Generating postgres instance...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("error generating postgres instance - %v\n", err)
		return
	}
	log.Println("PostgreSQL instance generated.")

	log.Println("Opening migration files...")
	fsrc, err := (&file.File{}).Open("file://database/migrations")
	if err != nil {
		log.Fatalf("error opening migration files - %v", err)
		return
	}
	log.Println("Migration files opened.")

	log.Println("Creating migration instance...")
	m, err := migrate.NewWithInstance("file", fsrc, "postgres", driver)
	if err != nil {
		log.Fatalf("error generating migrate instance - %v", err)
		return
	}
	log.Println("Migration instance created.")

	if *downFlag {
		log.Println("Rollback migration..")
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("error rollback migrations - %v", err)
			return
		}
		version, _, _ := m.Version()
		log.Printf("Rollback complete to version %d.\n", version)
	} else {
		log.Println("Migrating migration..")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("error migrating migrations - %v", err)
			return
		}
		log.Println("Migrate complete.")
	}
}
