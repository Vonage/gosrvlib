/*
Package uid provides functions to generate simple time-and-random-based unique identifiers.

It provides functions to generate 64 and 128 bit random identifiers in base-36 string format.

Deprecated: use github.com/Vonage/gosrvlib/pkg/uidc package instead.
*/
package uid

import "github.com/Vonage/gosrvlib/pkg/uidc"

// InitRandSeed initialize the random generator seed.
// You have to call this first and only once.
// Deprecated: use github.com/Vonage/gosrvlib/pkg/uidc package instead.
func InitRandSeed() error {
	return nil
}

// NewID64 generates and returns a new base-36-string-formatted 64 bit unique ID based on time (high 32 bit) and a random number (low 32 bit).
// NOTE: the zero time is set to the 1st of January of 10 years ago.
// Deprecated: use github.com/Vonage/gosrvlib/pkg/uidc package instead.
func NewID64() string {
	return uidc.NewID64()
}

// NewID128 generates and returns a new base-36-string-formatted 128 bit unique ID based on time (high 64 bit) and a random number (low 64 bit).
// Deprecated: use github.com/Vonage/gosrvlib/pkg/uidc package instead.
func NewID128() string {
	return uidc.NewID128()
}
