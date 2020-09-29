package gold

import (
	"bytes"
	"encoding/json"
	"regexp"
)

var (
	regexRFC3339    = regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z`)
	regexBcryptHash = regexp.MustCompile(`\$2[ayb]\$.{56}`)
	regexUUID       = regexp.MustCompile(`[0123456789abcdef-]{36}`)
)

// A Filter is a function that takes some content ([]byte) and returns a
// modified version of that content.
//
// This is useful for doing some work on the output before saving the content in
// a file.
type Filter func([]byte) []byte

// FilterFormatJSON formats the given JSON payload.
//
// It panics if it cannot parse the payload (invalid JSON).
func FilterFormatJSON(src []byte) []byte {
	dst := &bytes.Buffer{}

	err := json.Indent(dst, src, "", "\t")
	if err != nil {
		panic(err)
	}

	return dst.Bytes()
}

// FilterTimeRFC3339 swaps RFC3339 formatted strings with 0000-00-00T00:00:00Z.
func FilterTimeRFC3339(src []byte) []byte {
	return regexRFC3339.ReplaceAll(src, []byte("0000-00-00T00:00:00Z"))
}

// FilterBcryptHashes swaps bcrypt hashes with "_HASH_".
func FilterBcryptHashes(src []byte) []byte {
	return regexBcryptHash.ReplaceAll(src, []byte("_HASH_"))
}

// FilterUUIDs swaps UUIDs with "00000000-0000-0000-0000-000000000000".
func FilterUUIDs(src []byte) []byte {
	return regexUUID.ReplaceAll(src, []byte("00000000-0000-0000-0000-000000000000"))
}

// CustomFilter returns a filter that looks for patterns defined by reg and
// replaces them with rep.
func CustomFilter(reg *regexp.Regexp, rep string) Filter {
	return func(src []byte) []byte {
		return reg.ReplaceAll(src, []byte(rep))
	}
}
