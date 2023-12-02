// Package migrations contains code for launching postgresql migrations stored in the project
package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed postgres/*.sql
var embedMigrations embed.FS

func Migrate(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "postgres"); err != nil {
		panic(err)
	}
}
