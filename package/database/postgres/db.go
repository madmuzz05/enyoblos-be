package database

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MainDB struct {
	DB *gorm.DB
}

func Connect(host, user, password, dbname, port, sslmode string) (*MainDB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			&log.Logger, // gunakan zerolog logger
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set maximum number of idle connections
	sqlDB.SetMaxIdleConns(20)

	// Set maximum number of open connections
	sqlDB.SetMaxOpenConns(250)

	// Set the maximum lifetime of a connection (e.g., 10 minutes)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	return &MainDB{DB: db}, nil
}

func RunMigrationsPostgres(db *gorm.DB) {

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get sql.DB")
	}

	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create migration driver")
	}

	// Jalankan migration
	m, err := migrate.NewWithDatabaseInstance(
		"file://./package/database/postgres/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal().Err(err).Msg("migration init failed")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("migration failed")
	}

	// Migration selesai
	log.Info().Msg("Database migrations completed successfully")
}

// BeginTx memulai transaction dan return *gorm.DB
// Gunakan dengan defer untuk automatic rollback jika ada error
func (m *MainDB) BeginTx() *gorm.DB {
	tx := m.DB.Begin()
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("Failed to begin transaction")
		return nil
	}
	log.Info().Msg("Transaction started")
	return tx
}

// Commit melakukan commit transaction
func (m *MainDB) Commit(tx *gorm.DB) error {
	if tx == nil {
		return nil
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return err
	}
	log.Info().Msg("Transaction committed successfully")
	return nil
}

// Rollback melakukan rollback transaction
func (m *MainDB) Rollback(tx *gorm.DB) error {
	if tx == nil {
		return nil
	}
	if err := tx.Rollback().Error; err != nil {
		log.Error().Err(err).Msg("Failed to rollback transaction")
		return err
	}
	log.Info().Msg("Transaction rolled back")
	return nil
}

// WithTx menjalankan fungsi dalam transaction dengan defer
// Otomatis commit jika sukses, rollback jika error
func (m *MainDB) WithTx(fn func(*gorm.DB) error) error {
	tx := m.BeginTx()
	if tx == nil {
		return fmt.Errorf("failed to begin transaction")
	}

	// Defer untuk automatic rollback jika ada panic atau error
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("panic", r).Msg("Panic in transaction, rolling back")
			m.Rollback(tx)
			panic(r)
		}
	}()

	// Execute function
	if err := fn(tx); err != nil {
		log.Error().Err(err).Msg("Error in transaction, rolling back")
		m.Rollback(tx)
		return err
	}

	// Commit if no error
	return m.Commit(tx)
}

// TxCreate memulai transaction dan store dalam context
// Gunakan dengan defer TxSubmitTerr untuk automatic commit/rollback
func (m *MainDB) TxCreate(ctx *fiber.Ctx) *gorm.DB {
	tx := m.BeginTx()
	if tx == nil {
		log.Error().Msg("Failed to begin transaction")
		return nil
	}
	// Store transaction in context untuk diakses dalam defer
	ctx.Locals("tx", tx)
	return tx
}

// TxSubmitTerr melakukan commit atau rollback berdasarkan sysError
// Digunakan dalam defer statement setelah TxCreate
func (m *MainDB) TxSubmitTerr(ctx *fiber.Ctx, sysError syserror.SysError) {
	txInterface := ctx.Locals("tx")
	if txInterface == nil {
		return
	}

	tx, ok := txInterface.(*gorm.DB)
	if !ok {
		log.Error().Msg("Failed to retrieve transaction from context")
		return
	}

	if sysError != nil {
		// Rollback if there's an error
		m.Rollback(tx)
		log.Error().Msg("Transaction rolled back due to error")
	} else {
		// Commit if no error
		if err := m.Commit(tx); err != nil {
			log.Error().Err(err).Msg("Failed to commit transaction")
		}
	}
}
