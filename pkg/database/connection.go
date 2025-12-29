package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Database struct {
	Config *Config
	DSN    string
	DB     *gorm.DB
}

func NewDatabase(config Config) (*Database, error) {
	// set default connection values
	config.NormalizeValue()

	// generate dsn
	dsn, err := config.DSN()
	if err != nil {
		return nil, err
	}

	// set config

	db := Database{
		Config: &config,
		DSN:    dsn,
	}

	return &db, nil
}

// initial database
func (d *Database) Init() error {
	var level logger.LogLevel
	if d.Config.AppMode == "PRODUCTION" {
		level = logger.Error
	} else {
		level = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:        time.Second, // Slow SQL threshold
			LogLevel:             level,       // Log level
			ParameterizedQueries: true,        // Don't include params in the SQL log
			Colorful:             true,        // Disable color
		},
	)
	// create connection
	connection, err := gorm.Open(postgres.Open(d.DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	db, err := connection.DB()
	if err != nil {
		return err
	}

	// setting connection pool
	db.SetMaxIdleConns(*d.Config.MaxIdleConn)
	db.SetMaxOpenConns(*d.Config.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(*d.Config.MaxConnLifetime) * time.Second)

	d.DB = connection

	return nil
}

func (d *Database) Close() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
