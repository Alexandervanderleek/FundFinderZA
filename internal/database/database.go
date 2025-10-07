package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sqlx.DB
}

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDB(config *DbConfig) (*DB, error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)

	conn, err := sqlx.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %s", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error ping'ing database %s", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
