package main

import (
	"math/rand"
	"database/sql"
	"fmt"
)

func gen() (string, int) {
	return string('a' + uint8(rand.Uint32()%26)),
		int(rand.Uint32() % 3)
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
		panic(1)
	}
	fmt.Println("select", name, item)
	return i == 0
}
