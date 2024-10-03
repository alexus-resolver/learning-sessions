package main

func main() {
	panic_level2 := func() {
		panic("uh oh")
	}
	panic_level1 := func() {
		panic_level2()
	}

	panic_level1()
}

// go run 2_panics.go
