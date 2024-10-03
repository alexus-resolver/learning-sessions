1. why the error handling approach in GO is different? 
    - makes error processing routes explicit

2. panics
    - panics vs errors
        - panic is an unexpected error
        - panic indicates a bug in the code
    - panics recovery
        - defer and recover
    - panics in goroutine
    - panics in tests

3. ways of generating and consuming errors
    - error is an interface
    - ways to define error types
    - ignoring of any error is prohibited - log at least

4. rules of wrapping errors
    - as, is
    - to wrap or not to wrap
    - messages for wrapped errors

5. practices 
    - snippets
    - copilot

