package database

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	_dbCtxKey = "DBCtx"
)

var (
	_dbs       map[string]*sql.DB
	_defaultDB string
)

// Init initialize the database context internal state
// defaultDB is used to provide default database when database is not set like go routines where context may not be shared
func Init(dbs map[string]*sql.DB, defaultDB string) {
	_dbs = dbs
	_defaultDB = defaultDB
	if defaultDB == "" {
		panic("DefaultDB name is not correct" + defaultDB)
	}
	if len(dbs) == 0 {
		panic("DBs map can not be empty")
	}
	if _, ok := dbs[defaultDB]; !ok {
		panic("DefaultDB " + defaultDB + " is not set in the map")
	}
}

// FromContext this retries the database from context if it is not found than it returns default database
func FromContext(ctx context.Context) *sql.DB {
	db := ctx.Value(_dbCtxKey)
	if db == nil {
		return _dbs[_defaultDB]
	}
	return db.(*sql.DB)
}

// SetDb this is used to set the database context and same context should be passed across method calls to deal with database
func SetDb(ctx context.Context, name string) context.Context {
	db, found := _dbs[name]
	if !found {
		log.Println("DB name is not set using default one", name)
		db = _dbs[_defaultDB]
	}
	return context.WithValue(ctx, _dbCtxKey, db)
}
