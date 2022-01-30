package dbwrapper

import (
	"database/sql"
)

type SQLDB struct {
	DB *sql.DB
}

// func (db *SQLDB) Query(q string) {
// 	//debating whether to replace all models with functions that return strings,
//  //and just implement this
// }
