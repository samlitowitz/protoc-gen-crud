package test_cases

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	pgx_zap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

func PgSQLDBURLFromEnv() (string, error) {
	host := os.Getenv("DB_HOST")
	if len(host) == 0 {
		return "", fmt.Errorf("no host provided")
	}
	userFile := os.Getenv("DB_USER_FILE")
	if len(userFile) == 0 {
		return "", fmt.Errorf("no user file provided")
	}
	passwordFile := os.Getenv("DB_PASSWORD_FILE")
	if len(passwordFile) == 0 {
		return "", fmt.Errorf("no password file provided")
	}

	user, err := os.ReadFile(userFile)
	if err != nil {
		return "", fmt.Errorf("pgsql dburl from env: user: %w", err)
	}
	password, err := os.ReadFile(passwordFile)
	if err != nil {
		return "", fmt.Errorf("pgsql dburl from env: password: %w", err)
	}
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		strings.TrimSpace(string(user)),
		strings.TrimSpace(string(password)),
		host,
		strings.TrimSpace(string(user)),
	)

	connConfig, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return "", fmt.Errorf("parse config: %w", err)
	}
	devLogger, err := zap.NewDevelopment()
	if err != nil {
		return "", fmt.Errorf("new dev logger: %w", err)
	}
	connConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgx_zap.NewLogger(devLogger),
		LogLevel: tracelog.LogLevelError,
	}
	connStr := stdlib.RegisterConnConfig(connConfig)
	return connStr, nil
}

func PgSQLExecSQLFile(db *sql.DB, file string) error {
	code, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(code))
	if err != nil {
		return err
	}
	return nil
}
