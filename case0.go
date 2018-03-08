package main

import (
	"time"
	"fmt"
	"database/sql"
)

//会出现死锁
func case0 () {
	wg.Add(3)
	go run(1)
	go run(2)
	go run(3)
}

func run(id int) {
	for n := 100; n > 0; n-- {
		name, item := gen()
		now := time.Now()
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}
		if !query(tx, name, item) {
			fmt.Println(id, n, false, time.Since(now))
			tx.Commit()
			continue
		}
		r, err := tx.Exec("insert into user value (null,?,?)", name, item)
		if err != nil {
			panic(err)
		}
		a, err := r.RowsAffected()
		if err != nil {
			panic(err)
		}
		tx.Commit()
		fmt.Println(id, n, a == 1, time.Since(now))
	}
	wg.Done()
}

func query(tx *sql.Tx, name string, item int) bool {
	r, err := tx.Query("select count(*) from user where name=? and item=? for update", name, item)
	if err != nil {
		panic(err)
	}
	r.Next()
	i := 0
	r.Scan(&i)
	if r.Next() {
		fmt.Println("---", name, item)
		panic(1)
	}
	return i == 0
}
