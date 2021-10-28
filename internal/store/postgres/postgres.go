package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"ww/internal/conf"
	"ww/internal/store"
)

// Client is the database client
type Client struct {
	logger *zap.SugaredLogger
	db     *sqlx.DB
	newID  func() string
}

// New returns a new database client
func New(config *conf.Config) (*Client, error) {

	logger := zap.S().With("package", "store.postgres")

	var err error

	// Credentials
	var dbCreds string
	if username := config.Database.Username; username != "" {
		dbCreds += fmt.Sprintf("user=%s ", username)
	}
	if password := config.Database.Password; password != "" {
		dbCreds += fmt.Sprintf("password=%s ", password)
	}

	// Host + Port
	var connStr strings.Builder // Regular credentials
	if hostname := config.Database.Hostname; hostname != "" {
		connStr.WriteString(fmt.Sprintf("host=%s ", hostname))
	} else {
		return nil, fmt.Errorf("No hostname specified")
	}
	if port := strconv.Itoa(config.Database.Port); port != "" {
		connStr.WriteString(fmt.Sprintf("port=%s ", port))
	}

	// Database Name
	dbName := config.Database.Database

	var db *sqlx.DB

	// Connect to database
	connStr.WriteString(fmt.Sprintf("dbname=%s ", dbName))
	connConfig, err := pgx.ParseConfig(dbCreds + connStr.String())
	if err != nil {
		return nil, fmt.Errorf("could not parse pgx config: %s", err)
	}
	if config.Database.LogQueries {
		connConfig.Logger = &queryLogger{logger: logger}
	}

	for retries := config.Database.Retries; retries > 0; retries-- {

		// Attempt connecting to the database
		db, err = sqlx.Connect("pgx", stdlib.RegisterConnConfig(connConfig))
		if err == nil {
			// Ping the database
			if err = db.Ping(); err != nil {
				return nil, fmt.Errorf("could not ping database %w", err)
			}
			break // connected
		} else if strings.Contains(err.Error(), "connection refused") {
			logger.Warn("failed to connect to database, sleeping and retry")
			duration, _ := time.ParseDuration(config.Database.SleepBetweenRetries)
			time.Sleep(duration)
			continue
		}

		// Some other error
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("retries exausted, last error: %v", err)
	}

	db.SetMaxOpenConns(config.Database.MaxConnections)

	logger.Debugw("Connected to database server",
		"database.host", config.Database.Hostname,
		"database.username", config.Database.Username,
		"database.port", config.Database.Port,
		"database.database", config.Database.Database,
	)

	c := &Client{
		logger: logger,
		db:     db,
		newID: func() string {
			return xid.New().String()
		},
	}

	return c, nil

}

type queryLogger struct {
	logger *zap.SugaredLogger
}

func (ql *queryLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	ql.logger.Debugw(msg, "level", level, zap.Any("data", data))
}

// Lookup of postgres error codes to basic errors we can return to a user
var pgErrorCodeToStoreErrorType = map[string]store.ErrorType{
	"23502": store.ErrorTypeIncomplete,
	"23503": store.ErrorTypeForeignKey,
	"23505": store.ErrorTypeDuplicate,
	"23514": store.ErrorTypeInvalid,
}

func wrapError(err error) error {
	switch e := err.(type) {
	case *pgconn.PgError:
		if et, found := pgErrorCodeToStoreErrorType[e.Code]; found {
			return &store.Error{
				Type: et,
				Err:  err,
			}
		}
	}
	return err
}
