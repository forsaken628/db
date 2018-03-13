package main

import (
	"sync"
	"time"
	"fmt"
)

func case2() {
	c := sync.NewCond(&sync.Mutex{})

	wg.Add(n)
	es := make([]*Endpoint, n)
	for i := 0; i < n; i++ {
		es[i] = NewEndpoint(i, c, case2run)
	}

	for i := 0; i < n; i++ {
		name, item := gen()
		es[i].ch <- Msg{
			name: name,
			item: item,
		}
		<-es[i].done
	}

	c.Broadcast()
}

func case2run(id int, m Msg) {
	now := time.Now()
	//r, err := db.Exec("insert into user select null,?,? from dual where not exists(select 1 from user where name = ? and item = ?)", m.name, m.item, m.name, m.item)
	r, err := db.Exec("insert into user value (null,?,?)", m.name, m.item)
	if err != nil {
		fmt.Println(id, "insert", m.name, m.item)
		//panic(err)
		fmt.Println(err)
		wg.Done()
		return
	}
	a, err := r.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println(a == 1, time.Since(now))
	wg.Done()
}
