package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Begin is ...
func (db *DB) Begin() *DB {
	var d *gorm.DB
	depth := db.txDepth
	if depth == 0 {
		d = db.DB.Begin()
	} else {
		d = db.DB.Exec(fmt.Sprintf("SAVEPOINT LEVEL%d", depth))
	}

	return &DB{
		DB:      d,
		txDepth: depth + 1,
	}
}

// Commit is ...
func (db *DB) Commit() *DB {
	if db.txDepth == 0 {
		db.AddError(gorm.ErrInvalidTransaction)
		return db
	}

	depth := db.txDepth - 1
	if depth == 0 {
		db.DB.Commit()
	} else {
		db.DB.Exec(fmt.Sprintf("RELEASE SAVEPOINT LEVEL%d", depth))
	}
	return db
}

// Rollback is ...
func (db *DB) Rollback() *DB {
	if db.txDepth == 0 {
		db.AddError(gorm.ErrInvalidTransaction)
		return db
	}

	depth := db.txDepth - 1
	if depth == 0 {
		db.DB.Rollback()
	} else {
		db.DB.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT LEVEL%d", depth))
	}
	return db
}
