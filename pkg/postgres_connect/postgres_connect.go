package postgres_connect

import (
	"avito_test/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
	"log"
)

type PostgresStorage struct {
	Db *sql.DB
}

func NewPostgresStorage(cfg config.Postgres) (*PostgresStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	storage := &PostgresStorage{Db: db}

	if err := storage.runMigrations(cfg.MigrationPath); err != nil {
		return nil, fmt.Errorf("migrations failed: %v", err)
	}
	return storage, nil
}

func (s *PostgresStorage) runMigrations(path string) error {
	migrations := &migrate.FileMigrationSource{
		Dir: path,
	}

	ms := migrate.MigrationSet{
		TableName: "schema_migrations",
	}

	n, err := ms.Exec(s.Db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	if n == 0 {
		log.Println("No new migrations to apply.")
	} else {
		log.Printf("Applied %d database migration(s)\n", n)
	}
	return nil
}
