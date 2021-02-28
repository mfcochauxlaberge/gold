package gold_test

import (
	"regexp"
	"testing"

	. "github.com/mfcochauxlaberge/gold"

	"github.com/stretchr/testify/assert"
)

func TestFilterFormatJSON(t *testing.T) {
	assert := assert.New(t)

	src := []byte(`{"a": [0, "b"]}`)
	expected := `{
	"a": [
		0,
		"b"
	]
}`
	assert.Equal(expected, string(FilterFormatJSON(src)))

	// Invalid JSON will make FilterFormatJSON panic.
	assert.Panics(func() {
		_ = FilterFormatJSON([]byte(`this is invalid`))
	})
}

func TestFilterTimeRFC3339(t *testing.T) {
	assert := assert.New(t)

	src := []byte("The time is 2006-01-02T15:04:05Z07:00!")
	expected := "The time is 0000-00-00T00:00:00Z07:00!"
	assert.Equal(expected, string(FilterTimeRFC3339(src)))
}

func TestFilterBcryptHashes(t *testing.T) {
	assert := assert.New(t)

	hash := "$2a$10$aE6OdEySkK3g4HDnvJbFh.VXOW/gO7yDJEBqK/pnezwxjnkOo6kfC"

	src := []byte("The hash is " + string(hash) + "!")
	expected := "The hash is _HASH_!"
	assert.Equal(expected, string(FilterBcryptHashes(src)))
}

func TestFilterUUIDs(t *testing.T) {
	assert := assert.New(t)

	uuid := "25cef774-e2e6-4c78-ac1d-7b1acf1b28e6"

	src := []byte("The uuid is " + string(uuid) + "!")
	expected := "The uuid is 00000000-0000-0000-0000-000000000000!"
	assert.Equal(expected, string(FilterUUIDs(src)))
}

func TestCustomFilter(t *testing.T) {
	assert := assert.New(t)

	reg := regexp.MustCompile(`A[0-9]{1}Z`)
	filter := CustomFilter(reg, "__CUSTOM__")

	src := []byte("The string to be replaced is A4Z!")
	expected := "The string to be replaced is __CUSTOM__!"
	assert.Equal(expected, string(filter(src)))
}
