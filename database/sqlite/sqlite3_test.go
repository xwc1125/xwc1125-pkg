// description: sqlite
//
// @author: xwc1125
package sqlite

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/russross/meddler"
)

type Person struct {
	ID      int    `meddler:"id,pk"`
	Name    string `meddler:"name"`
	Age     int
	salary  int
	Created time.Time `meddler:"created,localtime"`
	Closed  time.Time `meddler:",localtimez"`
}

func TestDB(t *testing.T) {
	var db *sql.DB
	var err error
	// create the database
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic("error creating test database: " + err.Error())
	}

	person := &Person{
		Name: "Alice",
		Age:  22,
	}
	// 插入
	err = meddler.Insert(db, "person", person)
	if err != nil {
		panic(err)
	}
}
