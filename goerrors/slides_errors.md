# ERRORS
Error is one of two ways to handle errors in Go. <br>
The other one is panic.

Errors represent *expected* but not "positive" or "main flow" situations.

---


# CONVENTIONS
* An error is anything compatible with an error interface.
* Errors are usual values.
* Functions return an error as the last return argument.
* If the error is not nil, other return argument values should not be used.

```go
type error interface {
    Error() string
}

func CanReturnError() (any, error){...}
```

---

# ROOT ERROR
The app encountered an error situation for the first time.

```go
func Len(from,to int) (int,error){
    if from >= to {
        return "", errors.New("from should be less than to")
    }
    // ...
}
```

---

# TOOLS FOR ROOT ERRORS:

* errors.New("static error desc")
* fmt.Errorf("unsupported argument type %s", t)

---

# SENTINEL ERRORS

```go
// ErrNotAllEventsEnriched is returned when ...
var ErrNotAllEventsEnriched = errors.New("not all events ...")
return ErrNotAllEventsEnriched
```
* declares which errors to expect from the package
* allows to place documentation on what that error means, how to handle etc.

---

# ERROR HANDLING

```go
from, to, err := readBoundaries()
if err != nil{
    fmt.Println("Can not read boundaries: ", err)
    return // do not forget to return!
}
a, err := Len(from,to)
if err != nil {
    fmt.Println("Can not calc length: ", err)
    return
}
```

---

also can return the original error:

```go
a, err := Len(from,to)
if err != nil {
    return err
}
```

---

# ERROR HANDLING ANTI-PATTERNS:

---

# EXAMPLE 1: LOST ERROR VALUE

```go
file, err := os.ReadFile("log.txt")
if err != nil {
    return errors.New("read file error")
}
```

There is no way to know what actually happened.

---

# EXAMPLE 2: IGNORED ERROR

```go
err := db.Begin()
if err != nil { ... }

err := db.Exec(...)
if err != nil {
    db.Rollback()
    return err
}
```

Rollback returns a error too.

---

# EXAMPLE 3: MULTIPLE HANDLING
We should handle a error only once.
Logging - is handling.

```go
func A() error{
    err, data := json.Unmarshal(data, &target)
    if err != nil {
        log.Warn("can not unmarshal: ", err) // bad
        return errors.Wrap(err, "can not unmarshal")
    }
}
```

---

# ERROR HANDLING: WRAPPING

There are tools to produce a new error that will
contain another error plus additional description.

---

If the error can be handled at the current level, but still should stop the function's execution.

```go

func DoStuff() error{
    from,to := readFromTo()
    a, err := Len(from, to)
    if err != nil {
        return fmt.Errorf("calc len: %w", err)
        // return errors.Wrapf(err, "calc len for %i and %i", from, to) // legacy approach
    }
}
// Result: calc len for 2 and 1: from should be less than to
```

Use errors.Unwrap(err) to get the wrapped error.

---

Do not wrap-and-log errors to avoid it's multiple handling.

---

# CONVENTIONS FOR WRAPPING

---

# WHEN TO WRAP

* Simple rule: wrap always, it wouldn't hurt.

---

But
* Can skip wrapping if there is only one error-returning statement in the function
* Refactor the first return statement before adding a second one.

---

# WRAPPING MESSAGE 1/3

* wrapping statement describes the place in the current function
* not capitalized
* no punctuation at the end
* wrapped errors separated by ": " (for fmt.Errorf())

---

# WRAPPING MESSAGE 2/3

* avoid ''fail to'/'error at'/'can not' prefixes in wrapping messages

But ok for root errors and logging

```go
log.Warn("fail to read: %v", err)
errors.New("cant open socket")
```

---

# WRAPPING MESSAGE 3/3

a good practice:
wrapping statement is unique app-wide.

---

This is how it looks when handled:
`"fail to read form's data: get user: open db connection: network error 0x123"`

---

CHECKING WRAPPED ERRORS

---

```go
var ErrOops = errors.New("Oops")

func Foo() error {
    // do something...
    return ErrOops
}

func main() {
    err := Foo()
    if err == ErrOops {
        fmt.Println("All good, error handled")
    } else {
        panic(fmt.Errorf("unexpected error: %w" , err))
    }
}
```

All good, error handled

---

Situation 2: Foo was refactored.

```go
func Foo() error {
    return fmt.Errorf("error in foo: %w", ErrOops)
}
func main() {
    err := Foo()
    if err == ErrOops {
        fmt.Println("All good, error handled")
    } else {
        panic(fmt.Errorf("unexpected error: %w", err))
    }
}
```

`panic: unexpected error: error in foo: Oops`

---

Solution - use errors.Is()

```go
func main() {
    err := Foo()
    if errors.Is(err, ErrOops) {
        fmt.Println("All good, error handled")
    } else {
        panic(fmt.Errorf("unexpected error: %w", err))
    }
}
```

- Core still contains a few of these old-style, direct error checks.
- So, we have to check the function's usage when adding error wraps.

---

# ERRORS ARE NOT EXCEPTIONS

The 'exceptions' approach in many other languages mixes "unexpected" fatal errors and "flow control" errors.


In Go, they are separated into panics (for fatal cases) and errors (flow control).

---

That's why errors in their concept do not need any attached stack traces:
* They are not unexpected.
* An app is crafted to handle them already.

---


# UNUSUAL ERROR-HANDLING

* The usual error-handling case is a function that returns an error and we have a choice between wrapping and handling.
* Some functions can not return errors and this creates unusual error handling cases.

---

# EXAMPLES:
* main
* HTTP/grpc server handlers
* jobs of any scheduler/task/job manager
* the entry point of a background
* tests
* errors in loops

Such places require specific error handling code.

---

# RISKS:
* increased complexity
* duplication of error handling code
* forgotten returns

# RED FLAG:
Copy-pasted complex error handling code

---

# RECOMMENDATIONS:

* Reduce such places by extracting app logic to a function that returns an error.
* Document expected error handling procedures for new job managers of all kinds.
* Allow jobs and handlers to return errors in new systems.

---

# UNUSUAL ERROR-HANDLING:
# NESTED ERROR HANDLING

When using an error-returning function while handling another error

---

```go
err := foo()
if err != nil {
    undoErr := undo()
    if undoErr != nil {
        // return undoErr ???
    }
    // ...
}
```

* Which error to return?
    * They both can have wrapped sentinel errors that can be expected on the upper levels of the app.
* Do not lose error details here.
    * `errors.Join()`
    * `return fmt.Errorf("undo error '%s' after foo error: %w", undoErr, err)`


---

# FEEDBACK
* Should I make it more basic or more hardcore?
* Should I include some coding in a session?