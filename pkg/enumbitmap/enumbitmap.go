// Package enumbitmap provides functions to encode slices of enumeration strings into integer bitmap values and vice versa.
// Each bit correspond to a unique enumeration value.
//
// Example for uint8:
//
//    00000000 =   0 dec = NONE
//    00000001 =   1 dec = FIRST
//    00000010 =   2 dec = SECOND
//    00000100 =   4 dec = THIRD
//    00001000 =   8 dec = FOURTH
//    00010000 =  16 dec = FIFTH
//    00100000 =  32 dec = SIXTH
//    01000000 =  64 dec = SEVENTH
//    10000000 = 128 dec = EIGHT VALUE
//    00001001 = 1 + 8 = 9 dec = FIRST + FOURTH
package enumbitmap

import (
	"fmt"
)

const (
	// NumBitUint8 number of bit in uint8.
	NumBitUint8 = 8

	// NumBitUint16 number of bit in uint16.
	NumBitUint16 = 16

	// NumBitUint32 number of bit in uint32.
	NumBitUint32 = 32

	// NumBitUint64 number of bit in uint64.
	NumBitUint64 = 64
)

func mapIntToStrings(enum map[int]string, v uint64, numBit int) (s []string, err error) {
	if v == 0 {
		return []string{}, nil
	}

	s = make([]string, 0, numBit)
	errBits := make([]int, 0, numBit)

	var i uint64 = 1
	for bit := 1; bit <= numBit; bit++ {
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

// MapUint64ToStrings converts a uint64 bitmap value into a string slice.
func MapUint64ToStrings(enum map[int]string, v uint64) ([]string, error) {
	return mapIntToStrings(enum, v, NumBitUint64)
}

// MapUint32ToStrings converts a uint32 bitmap value into a string slice.
func MapUint32ToStrings(enum map[int]string, v uint32) ([]string, error) {
	return mapIntToStrings(enum, uint64(v), NumBitUint32)
}

// MapUint16ToStrings converts a uint16 bitmap value into a string slice.
func MapUint16ToStrings(enum map[int]string, v uint16) ([]string, error) {
	return mapIntToStrings(enum, uint64(v), NumBitUint16)
}

// MapUint8ToStrings converts a uint8 bitmap value into a string slice.
func MapUint8ToStrings(enum map[int]string, v uint8) ([]string, error) {
	return mapIntToStrings(enum, uint64(v), NumBitUint8)
}

func mapStringsToInt(enum map[string]int, s []string, numBit int) (v uint64, err error) {
	errStrings := make([]string, 0, numBit)

	for _, key := range s {
		id, ok := enum[key]
		if ok {
			v |= uint64(id)
		} else {
			errStrings = append(errStrings, key)
		}
	}

	if len(errStrings) > 0 {
		err = fmt.Errorf("unknown string values: %q", errStrings)
	}

	return v, err
}

// MapStringsToUint64 converts a string slice into a uint64 bitmap value.
func MapStringsToUint64(enum map[string]int, s []string) (uint64, error) {
	return mapStringsToInt(enum, s, NumBitUint64)
}

// MapStringsToUint32 converts a string slice into a uint32 bitmap value.
func MapStringsToUint32(enum map[string]int, s []string) (uint32, error) {
	v, err := mapStringsToInt(enum, s, NumBitUint32)
	return uint32(v), err
}

// MapStringsToUint16 converts a string slice into a uint16 bitmap value.
func MapStringsToUint16(enum map[string]int, s []string) (uint16, error) {
	v, err := mapStringsToInt(enum, s, NumBitUint16)
	return uint16(v), err
}

// MapStringsToUint8 converts a string slice into a uint8 bitmap value.
func MapStringsToUint8(enum map[string]int, s []string) (uint8, error) {
	v, err := mapStringsToInt(enum, s, NumBitUint8)
	return uint8(v), err
}
