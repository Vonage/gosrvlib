/*
Package phonekeypad provides functions to convert number strings to sequences of
numbers corresponding to a standard phone keypad.
*/
package phonekeypad

import (
	"fmt"
	"strings"
)

// KeypadDigit converts the input rune to a number corresponding to a standard phone keypad.
// If the input rune is not a number or a letter between A-Z or a-z, it returns -1 and false.
// Otherwise, it returns the corresponding number and true.
//
// The letter mapping is as follows:
//
//	+-----+-----+-----+
//	|  1  |  2  |  3  |
//	|     | ABC | DEF |
//	+-----+-----+-----+
//	|  4  |  5  |  6  |
//	| GHI | JKL | MNO |
//	+-----+-----+-----+
//	|  7  |  8  |  9  |
//	| PQRS| TUV | WXYZ|
//	+-----+-----+-----+
//	|     |  0  |     |
//	|     |     |     |
//	+-----+-----+-----+
func KeypadDigit(r rune) (int, bool) {
	if r >= '0' && r <= '9' {
		return int(r - '0'), true
	}

	if r >= 'a' && r <= 'z' {
		r -= ('a' - 'A')
	}

	return keypadAlphaToDigit(r)
}

// keypadAlphaToDigit converts the input alphabetical rune to a number corresponding to a standard phone keypad.
func keypadAlphaToDigit(r rune) (int, bool) {
	switch r {
	case 'A', 'B', 'C':
		return 2, true
	case 'D', 'E', 'F':
		return 3, true
	case 'G', 'H', 'I':
		return 4, true
	case 'J', 'K', 'L':
		return 5, true
	case 'M', 'N', 'O':
		return 6, true
	case 'P', 'Q', 'R', 'S':
		return 7, true
	case 'T', 'U', 'V':
		return 8, true
	case 'W', 'X', 'Y', 'Z':
		return 9, true
	}

	return -1, false
}

// KeypadNumber converts the input string to a sequence of numbers corresponding to a standard phone keypad.
// It skips any characters that are not numbers or letters between A-Z or a-z.
func KeypadNumber(num string) []int {
	seq := make([]int, 0, len(num))

	for _, r := range num {
		v, status := KeypadDigit(r)
		if status {
			seq = append(seq, v)
		}
	}

	return seq
}

// KeypadNumberString converts the input string to a sequence of numbers corresponding to a standard phone keypad.
// It skips any characters that are not numbers or letters between A-Z or a-z.
// It returns the sequence as a string.
func KeypadNumberString(num string) string {
	seq := KeypadNumber(num)

	return strings.Trim(
		strings.Join(
			strings.Fields(
				fmt.Sprint(seq),
			),
			"",
		),
		"[]",
	)
}
