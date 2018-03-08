package main

import "math/rand"

func gen() (string, int) {
	return string('a' + uint8(rand.Uint32()%26)),
		int(rand.Uint32() % 3)
}
