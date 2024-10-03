package main

import (
	"log"
	"runtime/debug"
	"testing"
	"time"
)

func TestPanic(t *testing.T) {

	panic_level2 := func() {
		panicking_function()
	}
	panic_level1 := func() {
		panic_level2()
	}

	panic_level1()

	// assert.Panics(t, func() {
	// 	panic_level1()
	// })
}

func panicking_function() {
	panic("uh oh")
}

func TestPanicIsBug(t *testing.T) {
	cases := []string{"a", "b", "a", "a", "b", "c"}
	// ...
	for i := 0; i < len(cases); i++ {
		switch cases[i] {
		case "a":
			log.Println("processing A")
		case "b":
			log.Println("processing B")
		default:
			panic("unexpected case")
		}
	}
}

func TestPanicRecovery(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			stacktrace := debug.Stack()
			//debug.PrintStack()
			log.Printf("Recovered! Panic:'%[1]v' \nPanic desc:'%#[1]v'\n", r)
			log.Printf("Stacktrace:\n%s", stacktrace)
		}
	}()

	panic(struct {
		msg    string
		reason string
	}{"uh oh", "something bad happened"})
}

func TestPanicInGoroutine(t *testing.T) {
	go func() {
		time.Sleep(5 * time.Second)
		panic("uh oh")
	}()

	for i := 0; i < 5; i++ {
		go func(worker int) {
			defer func() {
				log.Printf("worker %d still exits", worker)
			}()
			for {
				log.Printf("worker %d still alive", worker)
				time.Sleep(1 * time.Second)
			}
		}(i)
	}

	time.Sleep(10 * time.Second)

	// highlights:
	// ways to recover (in test, in goroutine)
	// research: how defers work in the case of panic
}
