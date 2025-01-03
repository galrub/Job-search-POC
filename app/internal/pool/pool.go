package pool

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/galrub/go/jobSearch/config"
	"github.com/galrub/go/jobSearch/internal/logger"
)

func DbConfig() *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	pd := config.GetPostgresData()

	dbConfig, err := pgxpool.ParseConfig(pd.CreateDSN())
	if err != nil {
		logger.LOG.Fatal().Err(err).Msg("Failed to create a pool config")
	}
	dbConfig.ConnConfig.User = pd.Username
	dbConfig.ConnConfig.Password = pd.Pw
	dbConfig.ConnConfig.Host = pd.Host
	dbConfig.ConnConfig.Database = pd.DbName
	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		//		logger.LOG.Debug().Msg("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		//		logger.LOG.Debug().Msg("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		//		logger.LOG.Debug().Msg("Closed the connection pool to the database!!")
	}

	return dbConfig
}

var DB *pgxpool.Pool = nil

func InitStatic() {
	pool, err := pgxpool.NewWithConfig(context.Background(), DbConfig())
	if err != nil {
		panic(err)
	}
	DB = pool
}

func GetConnection(ctx context.Context) (*pgxpool.Conn, error) {
	db, err := DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
