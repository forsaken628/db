package main

import "sync"

func case1() {
	mu:=&sync.Mutex{}

	c:=sync.NewCond(mu)

	c.Wait()
	c.Broadcast()

}
