package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/elgntt/effective-mobile/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func OpenDB(ctx context.Context, cfg config.DBConfig) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgPassword, cfg.PgDatabase,
	)

	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("OpenDB sqlx open: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("OpenDB sqlx ping: %w", err)
	}

	return db, nil
}

func Migrate(db *sqlx.DB, dbCfg config.DBConfig) error {
	driver, err := pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"pgx5", driver,
	)

	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
