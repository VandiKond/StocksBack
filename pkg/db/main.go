package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	ErrorCreateTable    = "error creating table"
	ErrorPreparingQuery = "error preparing query"
	ErrorInsertingUser  = "error inserting user"
	ErrorScanningRows   = "error scanning rows"
)

// The data base
type DB struct {
	db *sql.DB
}

// Creates table if not exists
func (db *DB) Create() error {
	query := `CREATE TABLE IF NOT EXISTS users (
   id BIGSERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   password VARCHAR(255) NOT NULL,
   solid_balance BIGINT DEFAULT 0,
   stock_balance BIGINT DEFAULT 0,
   is_blocked BOOLEAN DEFAULT FALSE,
   last_farming TIMESTAMP WITH TIME ZONE,
   created_at TIMESTAMP WITH TIME ZONE
);`
	_, err := db.db.Exec(query)
	if err != nil {
		return vanerrors.NewWrap(ErrorCreateTable, err, vanerrors.EmptyHandler)
	}
	return nil
}

// Creates a new user
func (db *DB) NewUser(u user_cfg.User) error {
	query := `INSERT INTO users (name, password, solid_balance, stock_balance, is_blocked, last_farming) VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Name, u.Password, u.SolidBalance, u.StockBalance, u.IsBlocked, u.LastFarming)
	if err != nil {
		return vanerrors.NewWrap(ErrorInsertingUser, err, vanerrors.EmptyHandler)
	}
	return nil
}

// Selecting
func (db *DB) Select(id uint64) (*user_cfg.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	var user user_cfg.User
	err = row.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, user.IsBlocked, &user.LastFarming, &user.CreatedAt)
	if err == nil {
		return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	return &user, nil
}

// Close the data base
func (db *DB) Close() error {
	return db.db.Close()
}
