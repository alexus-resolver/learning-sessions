package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefineError(t *testing.T) {

	var errorExample1 = errors.New("not enough arguments")
	assert.Error(t, errorExample1)
	var errorExample2 = fmt.Errorf("invalid size if the array (len is %d, expected %d)", 4, 5) // recommended (by me)
	assert.Error(t, errorExample2)
	// recommendation: do not split error message with variables. Place them in the end.
	// recommendation: reuse errors and make them public

	// examples at strconv: strconv.ErrRange

	// do not ignore errors:
	parsedBool, _ := strconv.ParseBool("true") // No!
	assert.True(t, parsedBool)

	parsedJson := 0
	json.Unmarshal([]byte(`42`), &parsedJson) // No!
	assert.Equal(t, 42, parsedJson)

	// the best way to handle errors in tests:
	var someError error = nil
	assert.NoError(t, someError)
	require.NoError(t, someError)
}
