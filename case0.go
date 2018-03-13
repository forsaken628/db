package main

import (
	"time"
	"fmt"
)

//空表会出现死锁
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
			fmt.Println("insert", name, item)
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
