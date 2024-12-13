package constructors

import (
	"github.com/vandi37/StocksBack/config/config"
	"github.com/vandi37/StocksBack/config/db_cfg"
	"github.com/vandi37/StocksBack/pkg/db"
	"github.com/vandi37/StocksBack/pkg/file_db"
	"github.com/vandi37/vanerrors"
)

// Errors
const (
	ConstructorNotFound = "constructor not found"
)

// It creates a new database
type Constructor interface {
	New(cfg config.DatabaseCfg, key string) (db_cfg.DataBase, error)
}

var db_constructors = map[string]Constructor{
	"postgres":    db.Constructor{},
	"postgressql": db.Constructor{},
	"file":        file_db.Constructor{},
	"fs":          file_db.Constructor{},
	"file system": file_db.Constructor{},
}

// Gets the constructor
func Get(s string) (Constructor, error) {
	c, ok := db_constructors[s]
	if !ok {
		return nil, vanerrors.NewSimple(ConstructorNotFound)
	}
	return c, nil
}
