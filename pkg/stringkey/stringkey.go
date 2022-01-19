// Package stringkey creates a key from multiple strings.
// NOTE: The key is based on a 64 bit hash and the input strings are trimmed and UTF8-normalized.
package stringkey

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	farmhash64 "github.com/tecnickcom/farmhash64/go/src"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	regexPatternEmptySpaces = `\s{1,}`
)

var (
	regexEmptySpaces = regexp.MustCompile(regexPatternEmptySpaces)
)

// StringKey stores the encoded key.
type StringKey struct {
	key uint64
}

// New encode (hash) input strings into a uint64 key.
func New(fields ...string) (*StringKey, error) {
	var b bytes.Buffer

	for _, v := range fields {
		b.WriteString(strings.ToLower(regexEmptySpaces.ReplaceAllLiteralString(strings.TrimSpace(v), " ")))
		b.WriteByte('\t') // separate input strings
	}

	utfTransformer := transform.Chain(norm.NFD, norm.NFC)

	nb, _, err := transform.Bytes(utfTransformer, b.Bytes())
	if err != nil {
		return nil, fmt.Errorf("unable to normalize the string: %w", err)
	}

	return &StringKey{key: farmhash64.FarmHash64(nb)}, nil
}

// Key returns a uint64 key.
func (sk *StringKey) Key() uint64 {
	return sk.key
}

// String returns a variable-length string key.
func (sk *StringKey) String() string {
	return strconv.FormatUint(sk.key, 36)
}

// Hex returns a fixed-length 16 digits hexadecimal string key.
func (sk *StringKey) Hex() string {
	s := strconv.FormatUint(sk.key, 16)

	slen := len(s)
	if slen < 16 {
		return strings.Repeat("0", (16-slen)) + s
	}

	return s
}
