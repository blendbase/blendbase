package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Logger *log.Logger
}

func NewApp() (*App, error) {
	// Set up the logger
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetLevel(log.DebugLevel) // TODO: Customize this depedening on the environment

	db, err := InitDb()
	if err != nil {
		log.Fatal("failed to connect database")
		return nil, err
	}

	return &App{DB: db, Logger: logger}, nil
}

func InitDb() (*gorm.DB, error) {
	log.Info("Initializing database")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
