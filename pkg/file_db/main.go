package file_db

import (
	"encoding/json"
	"io"
	"os"

	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/StocksBack/pkg/query"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	ErrorOpeningFile  = "error opening file"
	ErrorEncodingData = "error encoding data"
	ErrorDecodingData = "error decoding data"
	InvalidQuery      = "invalid query"
	InvalidId         = "invalid id"
)

// The file data base
type FileDB struct {
	*os.File
	data []user_cfg.User
	key  string
}

// Creates a new file data base
func NewFileDB(fn string, key string) (*FileDB, error) {
	// Opens the file
	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorOpeningFile, err, vanerrors.EmptyHandler)
	}

	return &FileDB{
		File: file,
		data: []user_cfg.User{},
		key:  key,
	}, nil
}

// Created tables (a array of users)
func (db *FileDB) Create() error {
	// Decoding data
	usrArr := []user_cfg.User{}
	err := json.NewDecoder(db).Decode(&usrArr)

	// setting data
	db.data = usrArr

	if err == io.EOF {
		// Saving data if the file is empty
		err = db.Save()
		if err != nil {
			return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
		}
	} else if err != nil {
		return vanerrors.NewWrap(ErrorDecodingData, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Saves the data in the file
func (db *FileDB) Save() error {
	// Marshals data
	jsonData, err := json.Marshal(db.data)
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}

	// Writes at 0 offset
	_, err = db.WriteAt(jsonData, 0)
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Created a new user
func (db *FileDB) NewUser(usr user_cfg.User) error {
	// Gets the user data
	usrArr := db.data

	// Checks the id
	if usr.Id != uint64(len(usrArr)) {
		return vanerrors.NewSimple(InvalidId)
	}

	// Ads the user
	usrArr = append(usrArr, usr)

	// Saving users to the data
	db.data = usrArr

	// Saving the data base
	err := db.Save()
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Gets all users
func (db *FileDB) GetAll() ([]user_cfg.User, error) {
	return db.data, nil
}

// Selecting user by id
func (db *FileDB) Select(id uint64) (*user_cfg.User, error) {
	usrArr := db.data
	return &usrArr[id], nil
}

// Selecting by query with limit
func (db *FileDB) SelectNumBy(q query.Query, num int) ([]user_cfg.User, error) {
	// Gets the user data
	usrArr := db.data

	// If num is les then zero sets it to length of users
	if num < 0 {
		num = len(usrArr)
	}

	// sorting bt query
	return q.Sort(usrArr, num)
}

// Selects users by query
func (db *FileDB) SelectBy(q query.Query) ([]user_cfg.User, error) {
	return db.SelectNumBy(q, -1)
}

// Selects user by query
func (db *FileDB) SelectOneBy(q query.Query) (*user_cfg.User, error) {
	// Selects one user
	res, err := db.SelectNumBy(q, 1)

	if err != nil {
		return nil, err
	}

	if len(res) < 1 {
		return nil, nil
	}

	return &res[0], nil
}

// Updates the user
func (db *FileDB) Update(usr user_cfg.User) error {
	// Gets the user data
	usrArr := db.data

	// Checks the if the id is valid
	if len(usrArr) <= int(usr.Id) {
		return vanerrors.NewSimple(InvalidId)
	}

	// Updates the user
	usrArr[usr.Id] = usr

	// Saving users to the data
	db.data = usrArr

	// Saving the data base
	err := db.Save()
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Updating a group of users
func (db *FileDB) UpdateGroup(users []user_cfg.User) error {
	// Gets the user data
	usrArr := db.data

	// Checks all users
	for _, usr := range users {
		// Checks the if the id is valid
		if len(usrArr) <= int(usr.Id) {
			return vanerrors.NewSimple(InvalidId)
		}
		// Updates the user
		usrArr[usr.Id] = usr
	}

	// Saving users to the data
	db.data = usrArr

	// Saving the data base
	err := db.Save()
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}

	return nil
}

// Gets the length of users
func (db *FileDB) GetLen() (uint64, error) {
	return uint64(len(db.data)), nil
}

func (db *FileDB) CheckKey(key string) (bool, error) {
	return db.key == key, nil
}
