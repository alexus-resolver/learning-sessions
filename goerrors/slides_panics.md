Two ways of error handling: panics and errors

# Panic: basics
* panic(any)
* Stops the function immediately, triggers defers 
* Then, the same for each caller
* Stops the app with a non-zero exit code

```go
func F1(){
    defer fmt.Println("f1 defer")
    fmt.Println("f1 start")
    panic("anything, usually a string or error") // line 8
    fmt.Println("f1 end")
}
func main(){
    defer fmt.Println("main defer")
    fmt.Println("main start")
    F1() // line 14
    fmt.Println("main end")
}
```

```
main start
f1 start
f1 defer
main defer
panic: anything, usually a string or error

goroutine 1 [running]:
main.F1()
        /.../main.go:8 +0xa8
main.main()
        /.../main.go:14 +0x98
exit status 2
```

# Panics: how do other goroutines stop?
When panic stops the app, other goroutine's defer blocks are not triggered.
```go
 func main() {
    fmt.Println("main start")
    defer fmt.Println("main defer")

    go Background()

    time.Sleep(time.Second * 3)
    panic("main panic")
}
func Background() {
    defer fmt.Println("Background defer")
    for {
        fmt.Println("Background")
        time.Sleep(time.Second)
    }
}
```

```
    main start
    Background
    Background
    Background
    main defer
    panic: main panic
```

# Panics: recover
Panics can be recovered.
It's only possible on the same goroutine.
```go
 func main() {
        fmt.Println("main start")
        defer fmt.Println("main defer before recover") // 7
        Foo()
        defer fmt.Println("main defer after panic") // 6

         fmt.Println("main end") // 5
}
func Foo() {
}
defer fmt.Println("foo defer before recover") // 4
defer func() { // 3
        err := recover()
        fmt.Println("foo defer recover: ", err)
}()
defer fmt.Println("foo defer after recover") // 2
panic("panic") // 1
defer fmt.Println("foo defer after panic")
```

```
main start
foo defer after recover
foo defer recover:  panic
foo defer before recover
main end
main defer after panic
main defer before recover
```

# Panics: When to use?

* It cannot be the primary way of error handling.
* Panic should mean that code changes are required (even if handled).
* Even handled panics signalling about a severe problem of code and can not be ignored in production.

## Cases:
* No config files and other startup checks 
* Unexpected data that was not got from IO or user.


# When is it OK to use recovery?

Main - to log all panics


In code that produces worker or handler goroutines, like web servers or schedulers - do not allow panic in one handler to crash the whole app.

# Panics from the runtime

There are a few specific cases when Go runtime generates panics:
* dividing by zero 
* type conversions 
* nil pointer
These are usual panics, and they can be handled.

# Non-recoverable errors
Some fatal errors are not recoverable. 
E.g.:
* Stack overflow, 
* Out of memory.
  