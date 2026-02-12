package postgress

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewRemoteDatabase(ctx context.Context, host string, port int, user string, password string, database string) (*PostgresDB, error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, database)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection with error: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database with error: %w", err)
	}

	db := PostgresDB{pool: pool}

	err = db.initTable(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init table in database with error: %w", err)
	}

	return &db, nil
}

func (db *PostgresDB) Close() {
	db.pool.Close()
}

// TODO: Transfer to userservice
func (db *PostgresDB) initTable(ctx context.Context) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction with error: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	createUserGropusTable := `CREATE TABLE IF NOT EXISTS user_groups (
	 user_group_id INT PRIMARY KEY
	 );`

	_, err = tx.Exec(ctx, createUserGropusTable)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to create user table: %w", err)
	}

	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
	 user_id VARCHAR(255) PRIMARY KEY,
	 user_group_id INT NOT NULL,
	 username TEXT UNIQUE NOT NULL,
	 full_name TEXT NOT NULL,
	 FOREIGN KEY (user_group_id) REFERENCES user_groups(user_group_id)
	 );`

	_, err = tx.Exec(ctx, createUsersTable)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to create user table: %w", err)
	}

	createAuthDataTable := `CREATE TABLE IF NOT EXISTS auth_data (
	 user_id VARCHAR(255) PRIMARY KEY,
	 username TEXT UNIQUE NOT NULL,
	 password_hash TEXT NOT NULL,
	 FOREIGN KEY (user_id) REFERENCES users(user_id)
	 );`

	_, err = tx.Exec(ctx, createAuthDataTable)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to create user table: %w", err)
	}

	createSessionsTable := `CREATE TABLE IF NOT EXISTS sessions (
	 user_id VARCHAR(255) NOT NULL,
	 session_id TEXT PRIMARY KEY,
	 valid_until TIMESTAMP NOT NULL,
	 session_state INT NOT NULL,
	 FOREIGN KEY (user_id) REFERENCES users(user_id)
	 );`

	_, err = tx.Exec(ctx, createSessionsTable)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to create user table: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction with error: %w", err)
	}

	return nil
}
