// Package postgres implements PostgreSQL connection using GORM
package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second

	SQLStateUniqueViolation     = "23505"
	SQLStateForeignKeyViolation = "23503"
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	DB *gorm.DB
}

type Config struct {
	User     string
	Password string
	DbName   string
	SslMode  bool
	Host     string
	Port     string
	LogEnv   string
}

func getDsn(cfg *Config) string {
	var sslMode string
	if cfg.SslMode {
		sslMode = "enabled"
	} else {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DbName,
		cfg.Port,
		sslMode,
	)
}

func newPg(opts ...Option) *Postgres {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	return pg
}

func getLogger(cfg *Config) logger.Interface {
	var logLevel logger.LogLevel

	switch cfg.LogEnv {
	case "local":
	case "dev":
		logLevel = logger.Info
	default:
		logLevel = logger.Error
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logLevel,
			Colorful: true,
		},
	)
}

// New -.
func New(cfg *Config, opts ...Option) (*Postgres, error) {
	dsn := getDsn(cfg)
	pg := newPg(opts...)
	l := getLogger(cfg)

	var (
		err error
		db  *gorm.DB
	)

	for pg.connAttempts > 0 {
		db, err = gorm.Open(gpostgres.Open(dsn), &gorm.Config{
			Logger: l, // Adjust log level as needed
		})
		if err == nil {
			// Configure connection pool
			sqlDB, err := db.DB()
			if err != nil {
				return nil, fmt.Errorf("postgres - New - db.DB(): %w", err)
			}

			sqlDB.SetMaxOpenConns(pg.maxPoolSize)
			sqlDB.SetMaxIdleConns(1) // Adjust as needed
			sqlDB.SetConnMaxLifetime(time.Hour)

			// Create a context for Ping
			ctx, cancel := context.WithTimeout(context.Background(), pg.connTimeout)

			// We perform ping and immediately cancel the context
			err = sqlDB.PingContext(ctx)

			cancel()

			if err == nil {
				pg.DB = db

				break
			}
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - New - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.DB != nil {
		sqlDB, err := p.DB.DB()
		if err != nil {
			panic(fmt.Errorf("postgres - Close - db.DB(): %w", err))
		}

		err = sqlDB.Close()
		if err != nil {
			panic(fmt.Errorf("postgres - Close - db.DB().Close: %w", err))
		}
	}
}
