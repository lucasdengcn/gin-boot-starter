package migrations

import (
	"errors"
	"gin001/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	// database driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Build DB schema
func Build() {
	cfg := config.GetConfig()
	m, err := migrate.New("file://migrations/schemas", cfg.DataSource.URL)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return
		}
		if errors.Is(err, migrate.ErrNilVersion) {
			return
		}
		log.Fatal(err)
	}
}
