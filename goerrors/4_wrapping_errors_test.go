package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"testing"

	goErrors "errors"

	errors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWrapErrors(t *testing.T) {

	err := ProcessFileExample()
	if err != nil {
		log.Println("error happen:", err)
	}

	// lets look to the all layers of this error
	var errUnwrap = err
	counter := 0
	for errUnwrap != nil {
		log.Printf("[%d] error: %s (%T)", counter, errUnwrap.Error(), errUnwrap)
		errUnwrap = errors.Unwrap(errUnwrap)
		counter++
	}

	// Best practices of wrapping:
	// 1. must wrap if function has more than one place where the error is returned
	// 2. wrapping message should not start with "error in...", "filed to...", "unable to..." etc.
	// 3. wrapping message must be unique in the function and describe a place in that function
	// 4. do not capitalize the first letter of the message
	// 5. do not add punctuation to the end of the message (wrap case) ore use colon and space (fmt.Errorf case)
	// 6. wrapping message should be in terms of the function, not in terms of the caller

	// now lets check type of the error
	// do not use "=="! Use errors.Is() or errors.As() instead.

	// in code
	if errors.Is(err, os.ErrNotExist) {
		log.Println("yes, this is our error")
	}

	var specificError *os.PathError
	if errors.As(err, &specificError) {
		log.Printf("yes, this is our error and there are details: Path=%v, Op=%v", specificError.Path, specificError.Op)
	}

	// in tests
	assert.True(t, errors.Is(err, os.ErrNotExist))
	assert.True(t, goErrors.Is(err, syscall.Errno(1))) // "no such file or directory" is 2
	assert.ErrorIs(t, err, syscall.Errno(1))

}

func ProcessFileExample() error {
	f, err := os.Open("non-existing-file.txt")
	if err != nil {
		return errors.Wrap(err, "open file") // Core's way
		// return fmt.Errorf("open file: %w", err) // a new way
	}

	// reading file to string
	_, err = io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read file contents: %w", err) // a new way
	}

	return nil
}
