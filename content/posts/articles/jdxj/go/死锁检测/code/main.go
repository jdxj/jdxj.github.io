package main

import (
	"time"

	"github.com/sasha-s/go-deadlock"
)

func main() {
	dl()
}

func dl() {
	m1 := &deadlock.Mutex{}
	m2 := &deadlock.Mutex{}

	go func() {
		lock(m2, m1)
	}()

	lock(m1, m2)
}

func lock(m1, m2 *deadlock.Mutex) {
	m1.Lock()
	time.Sleep(time.Second)

	m2.Lock()
	time.Sleep(time.Second)
	m2.Unlock()

	m1.Unlock()
}
