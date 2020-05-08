package tests

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/xgxw/foundation-go/database"
)

func MockDB(t *testing.T) (*database.DB, sqlmock.Sqlmock, func()) {
	var db *database.DB
	_db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open("mysql", _db)
	if err != nil {
		t.Fatalf("failed to connect to testing db: %s", err)
	}
	gormDB = gormDB.Debug()
	db = &database.DB{
		DB: gormDB,
	}
	teardown := func() {
		db.Close()
	}
	return db, mock, teardown
}
