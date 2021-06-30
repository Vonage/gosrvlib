// Package enumbitmap provides functions to encode slices of enumeration strings into integer bitmap values and vice versa.
// Each bit correspond to a unique enumeration value.
package enumbitmap

import (
	"fmt"

	"go.uber.org/multierr"
)

const (
	numBitUint8  = 8
	numBitUint16 = 16
	numBitUint32 = 32
	numBitUint64 = 64
)

// MapUint8ToStrings converts a uint8 bitmap value into a string slice.
func MapUint8ToStrings(enum map[int]string, v uint8) (s []string, err error) {
	if v == 0 {
		return []string{}, nil
	}

	s = make([]string, 0, numBitUint8)

	i := 1
	for bit := 1; bit <= numBitUint8; bit++ {
		if v&i == i {
			name, ok := enum[i]
			if ok {
				s = append(s, name)
			} else {
				err = multierr.append(err, fmt.Errorf("unknown bit: %d", bit))
			}
		}
		i = (i << 1)
	}

	return s, err
}

// MapStringsToUint8 converts a string slice into a uint8 bitmap value.
func MapStringsToUint8(enum map[string]int, s []string) (v uint8, err error) {
	for _, key := range s {
		id, ok := enum[key]
		if ok {
			v |= uint8(id)
		} else {
			err = multierr.append(err, fmt.Errorf("unknown enum: '%s'", key))
		}
	}

	return v, err
}
