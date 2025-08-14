package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(host, user, password, dbname, port, sslmode string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db

	// Auto-migrate
	// if err := DB.AutoMigrate(&voteEntity.Vote{}); err != nil {
	//      log.Println("AutoMigrate failed:", err)
	//      return err
	// }
	return DB, nil
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
