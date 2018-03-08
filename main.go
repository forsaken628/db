package main // import "github.com/forsaken628/db"

import (
	"math/rand"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"sync"
)

var wg sync.WaitGroup
var db *sql.DB

func init() {
	db0, err := sql.Open("mysql", "root@/test")
	if err != nil {
		panic(err)
	}
	err = db0.Ping()
	if err != nil {
		panic(err)
	}
	db = db0
	rand.Seed(time.Now().UnixNano())
}

func main() {
	case0()
	wg.Wait()
}
