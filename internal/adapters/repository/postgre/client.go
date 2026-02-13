package postgre

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-auth/internal/adapters/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var ErrInvalidDSN = errors.New("invalid DSN provided")

type Client struct {
	DB *gorm.DB
}

func (c *Client) Ping(ctx context.Context) error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

func NewPostgreSQLClient(cfg *config.DBConfig) (*Client, error) {
	if cfg.DSN() == "" {
		return nil, ErrInvalidDSN
	}

	db, err := openPostgreSQLDB(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{DB: db}, nil
}

func openPostgreSQLDB(cfg *config.DBConfig) (*gorm.DB, error) {
	dialector := postgres.Open(cfg.DSN())

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             cfg.SlowThreshold(),
			LogLevel:                  cfg.LogLevel(),
			IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError(),
			Colorful:                  cfg.Colorful(),
		},
	)

	// GORM configuration
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         gormLogger,
	}

	// Open database connection
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}

	// Ensure database conenction is reachable
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database is unreachable: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns())
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns())
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime())
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime())
	return db, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB for closing: %w", err)
	}
	return sqlDB.Close()
}
