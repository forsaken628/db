package main

import (
	"sync"
	"fmt"
	"time"
)

func case1() {
	c := sync.NewCond(&sync.Mutex{})

	wg.Add(n)
	es := make([]*Endpoint, n)
	for i := 0; i < n; i++ {
		es[i] = NewEndpoint(i, c, case1run)
	}

	name, item := gen()
	for i := 0; i < n; i++ {
		es[i].ch <- Msg{
			name: name,
			item: item,
		}
		<-es[i].done
	}

	c.Broadcast()
}

func case1run(id int, m Msg) {
	now := time.Now()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	if !query(tx, m.name, m.item) {
		fmt.Println(id, false, time.Since(now))
		tx.Rollback()
		wg.Done()
		return
	}
	r, err := tx.Exec("insert into user value (null,?,?)", m.name, m.item)
	if err != nil {
		fmt.Println(id, "insert", m.name, m.item)
		panic(err)
	}
	a, err := r.RowsAffected()
	if err != nil {
		panic(err)
	}
	tx.Commit()
	fmt.Println(a == 1, time.Since(now))
	wg.Done()
}

type Msg struct {
	name string
	item int
}

type Endpoint struct {
	id   int
	c    *sync.Cond
	ch   chan Msg
	done chan struct{}
	fn   func(int, Msg)
}

func NewEndpoint(id int, c *sync.Cond, fn func(int, Msg)) *Endpoint {
	ch := make(chan Msg, 1)
	done := make(chan struct{})
	e := &Endpoint{
		id:   id,
		c:    c,
		ch:   ch,
		done: done,
		fn:   fn,
	}
	go e.loop()
	return e
}

func (e *Endpoint) loop() {
	for {
		d := <-e.ch
		e.c.L.Lock()
		e.done <- struct{}{}
		e.c.Wait()
		e.c.L.Unlock()
		e.fn(e.id, d)
	}
}
