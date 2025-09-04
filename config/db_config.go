package config

import (
	"fmt"
	"task-manager-app/constants"
	"task-manager-app/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		ApplicationConfig.PostgresHost,
		ApplicationConfig.Username,
		ApplicationConfig.Password,
		ApplicationConfig.DbName,
		ApplicationConfig.PostgresPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		utils.Sugar.Fatal(constants.ErrFailedToConnectDB+":", err)
	}

	// Set connection pool
	sqlDB, err := db.DB()
	if err != nil {
		utils.Sugar.Fatal(constants.ErrFailedToGetSqlDB+":", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	DB = db
	utils.Sugar.Info("Connected to PostgreSQL successfully")
}
