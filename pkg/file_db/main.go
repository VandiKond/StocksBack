package file_db

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/VandiKond/StocksBack/config/user_cfg"
	"github.com/VandiKond/vanerrors"
)

// The errors
const (
	ErrorOpeningFile  = "error opening file"  // error opening file
	ErrorCreatingData = "error creating data" // error creating data
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
	if err != nil {
		err := json.NewEncoder(db).Encode(usrArr)
		if err != nil {
			return vanerrors.NewWrap(ErrorCreatingData, err, vanerrors.EmptyHandler)
		}
	}
	return nil
}

func (db *FileDB) NewUser(usr user_cfg.User) error {
	usrArr := db.data

	usrArr = append(usrArr, usr)
	err := json.NewEncoder(db).Encode(usrArr)
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}
	db.data = usrArr
	return nil
}

func (db *FileDB) GetAll() ([]user_cfg.User, error) {
	return db.data, nil
}

func (db *FileDB) Select(uid uint64) (*user_cfg.User, error) {
	usrArr := db.data
	return &usrArr[uid], nil
}

func (db *FileDB) SelectNumBy(q any, num int) ([]user_cfg.User, error) {
	usrArr := db.data
	if num == -1 {
		num = len(usrArr) - 1
	}
	query, ok := q.([]QuerySetting)
	if !ok {
		return nil, vanerrors.NewSimple(InvalidQuery, fmt.Sprintf("expected query format %T, but go  %T", query, q))
	}
	var res []user_cfg.User

	for _, u := range usrArr {
		var do_append bool
		var cur_separator int
		for _, qr := range query {
			switch qr.Separator {
			case NOT_SEPARATOR:
				var is_true bool
				switch qr.Type {
				case ID:
					is_true = qr.Run(u.Id)
				case NAME:
					is_true = qr.Run(u.Name)
				case PASSWORD:
					is_true = qr.Run(u.Password)
				case SOLID_BALANCE:
					is_true = qr.Run(u.SolidBalance)
				case STOCK_BALANCE:
					is_true = qr.Run(u.StockBalance)
				case IS_BLOCKED:
					is_true = qr.Run(u.IsBlocked)
				case CREATED_AT:
					is_true = qr.Run(u.CreatedAt)
				default:
					return nil, vanerrors.NewSimple(InvalidQuery, "invalid type")
				}
				switch cur_separator {
				case OR:
					do_append = do_append || is_true
				case AND:
					do_append = do_append && is_true
				default:
					return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
				}
			case OR, AND:
				if cur_separator != NOT_SEPARATOR {
					return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
				}
			}
			cur_separator = qr.Separator
		}
		if cur_separator != NOT_SEPARATOR {
			return nil, vanerrors.NewSimple(InvalidQuery, "invalid order")
		}
		if do_append {
			res = append(res, u)
		}

		if len(res) > num {
			break
		}
	}
	return res, nil
}

func (db *FileDB) SelectBy(q any) ([]user_cfg.User, error) {
	return db.SelectNumBy(q, -1)
}

func (db *FileDB) SelectOneBy(q any) (*user_cfg.User, error) {
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

	err := json.NewEncoder(db).Encode(usrArr)
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}
	db.data = usrArr
	return nil
}

func (db *FileDB) UpdateGroup(users []user_cfg.User) error {
	usrArr := db.data
	for _, usr := range users {
		if len(usrArr) >= int(usr.Id) {
			return vanerrors.NewSimple(InvalidId)
		}
		usrArr[usr.Id] = usr
	}

	err := json.NewEncoder(db).Encode(usrArr)
	if err != nil {
		return vanerrors.NewWrap(ErrorEncodingData, err, vanerrors.EmptyHandler)
	}
	db.data = usrArr
	return nil
}

func (db *FileDB) GetLast() (uint64, error) {
	return db.data[len(db.data)-1].Id, nil
}
