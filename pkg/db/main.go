package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/vandi37/StocksBack/config/config"
	"github.com/vandi37/StocksBack/config/db_cfg"
	"github.com/vandi37/StocksBack/config/user_cfg"
	"github.com/vandi37/StocksBack/pkg/query"
	"github.com/vandi37/vanerrors"
)

// The errors
const (
	ErrorOpeningDataBase = "error opining database"
	ErrorCreateTable     = "error creating table"
	ErrorPreparingQuery  = "error preparing query"
	ErrorInsertingUser   = "error inserting user"
	ErrorScanningRows    = "error scanning rows"
	ErrorSelecting       = "error selecting"
	ErrorGettingLength   = "error getting length"
	ErrorUpdatingUser    = "error updating user"
	NotFound             = "not found"
)

// The data base
type DB struct {
	db  *sql.DB
	key string
}

// The db constructor
type Constructor struct{}

// Creates a new data base connection
func (c Constructor) New(cfg config.DatabaseCfg, key string) (db_cfg.DataBase, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name))
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorOpeningDataBase, err, vanerrors.EmptyHandler)
	}
	return &DB{db: db, key: key}, nil
}

// Creates table if not exists
func (db *DB) Init() error {

	query := `CREATE TABLE IF NOT EXISTS users (
		id BIGINT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		solid_balance BIGINT DEFAULT 0,
		stock_balance BIGINT DEFAULT 0,
		is_blocked BOOLEAN DEFAULT FALSE,
		last_farming TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.db.Exec(query)
	if err != nil {
		return vanerrors.NewWrap(ErrorCreateTable, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Creates a new user
func (db *DB) Create(u user_cfg.User) error {
	// Prepares the query
	query := `insert into users (id, name, password, solid_balance, stock_balance, is_blocked, last_farming, created_at) values ($1, $2, $3, $4, $5, $6, $7, $8);`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}

	defer stmt.Close()

	// Creates user
	_, err = stmt.Exec(u.Id, u.Name, u.Password, u.SolidBalance, u.StockBalance, u.IsBlocked, u.LastFarming, u.CreatedAt)
	if err != nil {
		return vanerrors.NewWrap(ErrorInsertingUser, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Gets all
func (db *DB) GetAll() ([]user_cfg.User, error) {
	query := `select * from users;`
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	defer rows.Close()

	var res []user_cfg.User

	// Adding users
	for rows.Next() {

		var user user_cfg.User
		err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, &user.IsBlocked, &user.LastFarming, &user.CreatedAt)
		if err != nil {
			return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
		}

		res = append(res, user)
	}

	if rows.Err() != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	return res, nil
}

// Selects all by query
func (db *DB) GetAllBy(q query.Query) ([]user_cfg.User, error) {
	return db.GetNumBy(q, -1)
}

// Selecting by query
func (db *DB) GetNumBy(q query.Query, num int) ([]user_cfg.User, error) {
	// Getting query part
	str, args := q.PrepareString()
	query := `select * from users where ` + str

	// Checking is the limit need
	if num > 0 {
		str += " limit " + strconv.Itoa(num)
	}

	str += ";"

	// Preparing query
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()

	// Getting rows
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	defer rows.Close()

	var res []user_cfg.User

	// Adding users
	for rows.Next() {

		var user user_cfg.User
		err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, &user.IsBlocked, &user.LastFarming, &user.CreatedAt)
		if err != nil {
			return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
		}

		res = append(res, user)
	}

	if rows.Err() != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	return res, nil
}

// Selecting
func (db *DB) GetOne(id uint64) (*user_cfg.User, error) {
	// Prepares the query
	query := `select * from users where id = $1;`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()

	// Selects the user
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	defer rows.Close()

	// Getting user
	var user user_cfg.User

	if !rows.Next() {
		return nil, vanerrors.NewSimple(NotFound)
	}

	err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, &user.IsBlocked, &user.LastFarming, &user.CreatedAt)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	if rows.Err() != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	return &user, nil
}

// Selects user by query
func (db *DB) GetOneBy(q query.Query) (*user_cfg.User, error) {
	// Selects one user
	res, err := db.GetNumBy(q, 1)

	if err != nil {
		return nil, err
	}

	if len(res) < 1 {
		return nil, nil
	}

	return &res[0], nil
}

// Updates block
func (db *DB) UpdateBlock(id uint64, block bool) (*user_cfg.User, error) {
	query := `update users set is_blocked = $1 where id = $2 returning *;`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()

	// Updating the user
	rows, err := stmt.Query(block, id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	defer rows.Close()

	// Getting user
	var user user_cfg.User

	if !rows.Next() {
		return nil, vanerrors.NewSimple(NotFound)
	}

	err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, &user.IsBlocked, &user.LastFarming, &user.CreatedAt)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	if rows.Err() != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	return &user, nil
}

// Updates last farm
func (db *DB) UpdateLastFarm(id uint64) (*user_cfg.User, error) {
	query := `update users set last_farming = $1 where id = $2 returning *;`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()

	// Updating the user
	rows, err := stmt.Query(time.Now(), id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	defer rows.Close()

	// Getting user
	var user user_cfg.User

	if !rows.Next() {
		return nil, vanerrors.NewSimple(NotFound)
	}

	err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, &user.IsBlocked, &user.LastFarming, &user.CreatedAt)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	if rows.Err() != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	return &user, nil
}

// Updates name
func (db *DB) UpdateName(id uint64, name string) (*user_cfg.User, error) {
	query := `update users set name = $1 where id = $2 returning *;`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()

	// Updating the user
	rows, err := stmt.Query(name, id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	defer rows.Close()

	// Getting user
	var user user_cfg.User

	if !rows.Next() {
		return nil, vanerrors.NewSimple(NotFound)
	}

	err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, &user.IsBlocked, &user.LastFarming, &user.CreatedAt)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	if rows.Err() != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	return &user, nil
}

// Updates password
func (db *DB) UpdatePassword(id uint64, password string) (*user_cfg.User, error) {
	query := `update users set password = $1 where id = $2 returning *;`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()

	// Updating the user
	rows, err := stmt.Query(password, id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	defer rows.Close()

	// Getting user
	var user user_cfg.User

	if !rows.Next() {
		return nil, vanerrors.NewSimple(NotFound)
	}

	err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, &user.IsBlocked, &user.LastFarming, &user.CreatedAt)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	if rows.Err() != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	return &user, nil
}

func (db *DB) UpdateSolids(id uint64, num int64) (*user_cfg.User, error) {
	query := "update users set solid_balance = solid_balance + $1 where id = $2 returning *;"

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()

	// Updating the user
	rows, err := stmt.Query(num, id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	defer rows.Close()

	// Getting user
	var user user_cfg.User

	if !rows.Next() {
		return nil, vanerrors.NewSimple(NotFound)
	}

	err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, &user.IsBlocked, &user.LastFarming, &user.CreatedAt)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	if rows.Err() != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	return &user, nil
}

func (db *DB) UpdateStocks(id uint64, num int64) (*user_cfg.User, error) {
	query := "update users set stock_balance = stock_balance + $1 where id = $2 returning *;"

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorPreparingQuery, err, vanerrors.EmptyHandler)
	}
	defer stmt.Close()

	// Updating the user
	rows, err := stmt.Query(num, id)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	defer rows.Close()

	// Getting user
	var user user_cfg.User

	if !rows.Next() {
		return nil, vanerrors.NewSimple(NotFound)
	}

	err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.SolidBalance, &user.StockBalance, &user.IsBlocked, &user.LastFarming, &user.CreatedAt)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorScanningRows, err, vanerrors.EmptyHandler)
	}

	if rows.Err() != nil {
		return nil, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}

	return &user, nil
}

func (db *DB) Len() (uint64, error) {
	query := `select id from users order by id desc limit 1;`
	// Updating the user
	rows, err := db.db.Query(query)
	if err != nil {
		return 0, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	defer rows.Close()

	var length uint64

	if !rows.Next() {
		return 0, nil
	}

	err = rows.Scan(&length)

	if err != nil {
		return 0, vanerrors.NewWrap(ErrorSelecting, err, vanerrors.EmptyHandler)
	}
	return length + 1, nil
}

// Close the data base
func (db *DB) Close() error {
	return db.db.Close()
}

// Checks key
func (db *DB) CheckKey(key string) (bool, error) {
	return db.key == key, nil
}
