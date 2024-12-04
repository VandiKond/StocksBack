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
	ErrorOpeningFile  = "error opening file"  // error opening file
	ErrorEncodingData = "error encoding data" // error encoding data
	ErrorDecodingData = "error decoding data" // error encoding data
	InvalidQuery      = "invalid query"       // invalid query
	InvalidId         = "invalid id"          // invalid id
)

type FileDB struct {
	*os.File
	data []user_cfg.User
}

func NewFileDB(fn string) (*FileDB, error) {
	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorOpeningFile, err, vanerrors.EmptyHandler)
	}
	return &FileDB{
		File: file,
		data: []user_cfg.User{},
	}, nil
}

func (db *FileDB) Create() error {
	usrArr := []user_cfg.User{}
	err := json.NewDecoder(db).Decode(&usrArr)
	db.data = usrArr
	if err == io.EOF {
		err = db.Save()
		if err != nil {
			return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
		}
	} else if err != nil {
		return vanerrors.NewWrap(ErrorDecodingData, err, vanerrors.EmptyHandler)
	}
	return nil
}

func (db *FileDB) Save() error {
	jsonData, err := json.Marshal(db.data)
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}
	_, err = db.WriteAt(jsonData, 0)
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}
	return nil
}

func (db *FileDB) NewUser(usr user_cfg.User) error {
	usrArr := db.data

	if usr.Id != uint64(len(usrArr)) {
		return vanerrors.NewSimple(InvalidId)
	}
	usrArr = append(usrArr, usr)
	db.data = usrArr
	err := db.Save()
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}
	return nil
}

func (db *FileDB) GetAll() ([]user_cfg.User, error) {
	return db.data, nil
}

func (db *FileDB) Select(uid uint64) (*user_cfg.User, error) {
	usrArr := db.data
	return &usrArr[uid], nil
}

func (db *FileDB) SelectNumBy(q query.Query, num int) ([]user_cfg.User, error) {
	usrArr := db.data
	if num == -1 {
	}
	return q.Sort(usrArr, num)

}

func (db *FileDB) SelectBy(q query.Query) ([]user_cfg.User, error) {
	return db.SelectNumBy(q, -1)
}

func (db *FileDB) SelectOneBy(q query.Query) (*user_cfg.User, error) {
	res, err := db.SelectNumBy(q, 1)
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, nil
	}

	return &res[0], nil
}

func (db *FileDB) Update(usr user_cfg.User) error {
	usrArr := db.data

	if len(usrArr) >= int(usr.Id) {
		return vanerrors.NewSimple(InvalidId)
	}
	usrArr[usr.Id] = usr

	db.data = usrArr
	err := db.Save()
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}
	return nil
}

func (db *FileDB) UpdateGroup(users []user_cfg.User) error {
	usrArr := db.data
	for _, usr := range users {
		if len(usrArr) < int(usr.Id) {
			return vanerrors.NewSimple(InvalidId)
		}
		usrArr[usr.Id] = usr
	}
	db.data = usrArr
	err := db.Save()
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}
	return nil
}

func (db *FileDB) GetLen() (uint64, error) {
	return uint64(len(db.data)), nil
}
