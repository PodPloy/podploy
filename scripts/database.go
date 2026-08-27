package scripts

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5"
)

type databaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	AdminDB  string
	SSLMode  string
}

func (c databaseConfig) connectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.AdminDB, c.SSLMode,
	)
}

type databaseManager struct {
	conn *pgx.Conn
}

func newDatabaseManager(ctx context.Context, config databaseConfig) (*databaseManager, error) {
	conn, err := pgx.Connect(ctx, config.connectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to establish database connection: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &databaseManager{conn: conn}, nil
}

func (m *databaseManager) close(ctx context.Context) error {
	return m.conn.Close(ctx)
}

func (m *databaseManager) DatabaseExists(ctx context.Context, dbname string) (bool, error) {
	const query = "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1);"

	var exists bool

	if err := m.conn.QueryRow(ctx, query, dbname).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to query database existence: %w", err)
	}

	return exists, nil
}

func validateDBName(name string) error {
	validPattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]{0,62}$`)

	if !validPattern.MatchString(name) {
		return fmt.Errorf("invalid database name '%s': must start with a letter/underscore and container only alphanumeric characters", name)
	}

	return nil
}
