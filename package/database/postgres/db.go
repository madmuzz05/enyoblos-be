package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/source/file" // WAJIB
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/rs/zerolog/log"
)

type MainDB struct {
	DB *sqlx.DB
}

type DBWithCtx struct {
	DB  *sqlx.DB
	Ctx context.Context
}

func Connect(host, user, password, dbname, port, sslmode string) (*MainDB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Pool config
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(250)
	db.SetConnMaxLifetime(10 * time.Minute)

	return &MainDB{DB: db}, nil
}
func (m *MainDB) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := m.DB.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil, err
	}
	return tx, nil
}

func (m *MainDB) Commit(tx *sqlx.Tx) error {
	if tx == nil {
		return nil
	}
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return err
	}
	return nil
}

func (m *MainDB) Rollback(tx *sqlx.Tx) error {
	if tx == nil {
		return nil
	}
	if err := tx.Rollback(); err != nil {
		log.Error().Err(err).Msg("Failed to rollback transaction")
		return err
	}
	return nil
}

func (m *MainDB) WithTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := m.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("panic", r).Msg("panic in tx")
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func TxCreate(c fiber.Ctx, db *sqlx.DB) (*sqlx.Tx, syserror.SysError) {
	tx, err := db.BeginTxx(c.Context(), nil)
	if err != nil {
		return nil, syserror.CreateError(err, fiber.StatusInternalServerError, "Gagal membuat transaction")
	}

	c.Locals("tx", tx)
	return tx, nil
}

func TxSubmitTerr(c fiber.Ctx, sysError syserror.SysError) {
	txInterface := c.Locals("tx")
	if txInterface == nil {
		return
	}

	tx, ok := txInterface.(*sqlx.Tx)
	if !ok {
		log.Error().Msg("invalid tx type")
		return
	}

	if sysError != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err).Msg("failed rollback transaction")
		} else {
			log.Info().Msg("transaction rolled back")
		}
	} else {
		if err := tx.Commit(); err != nil {
			log.Error().Err(err).Msg("failed commit transaction")
		} else {
			log.Info().Msg("transaction committed")
		}
	}
}

func RunMigrationsPostgres(db *sqlx.DB) {

	sqlDB := db.DB

	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create migration driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./package/database/postgres/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("migration init failed")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("migration failed")
	}

	log.Info().Msg("Database migrations completed successfully")
}

func (db *DBWithCtx) Get(dest interface{}, query string, args ...interface{}) error {
	return db.DB.GetContext(db.Ctx, dest, query, args...)
}

func (db *DBWithCtx) Select(dest interface{}, query string, args ...interface{}) error {
	return db.DB.SelectContext(db.Ctx, dest, query, args...)
}

func (db *DBWithCtx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.DB.ExecContext(db.Ctx, query, args...)
}
