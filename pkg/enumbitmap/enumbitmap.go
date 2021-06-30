// Package enumbitmap provides functions to encode slices of enumeration strings into integer bitmap values and vice versa.
// Each bit correspond to a unique enumeration value.
package enumbitmap

import (
	"fmt"
)

const (
	numBitUint8 = 8
)

// MapUint8ToStrings converts a uint8 bitmap value into a string slice.
func MapUint8ToStrings(enum map[int]string, v uint8) (s []string, err error) {
	if v == 0 {
		return []string{}, nil
	}

	s = make([]string, 0, numBitUint8)
	errBits := make([]int, 0, numBitUint8)

	var i uint8 = 1
	for bit := 1; bit <= numBitUint8; bit++ {
		if v&i == i {
			name, ok := enum[int(i)]
			if ok {
				s = append(s, name)
			} else {
				errBits = append(errBits, bit)
			}
		}

		i = (i << 1)
	}

	if len(errBits) > 0 {
		err = fmt.Errorf("unknown bit values: %v", errBits)
	}

	return s, err
}

// MapStringsToUint8 converts a string slice into a uint8 bitmap value.
func MapStringsToUint8(enum map[string]int, s []string) (v uint8, err error) {
	errStrings := make([]string, 0, numBitUint8)

	for _, key := range s {
		id, ok := enum[key]
		if ok {
			v |= uint8(id)
		} else {
			errStrings = append(errStrings, key)
		}
	}

	if len(errStrings) > 0 {
		err = fmt.Errorf("unknown string values: %v", errStrings)
	}

	return v, err
}
