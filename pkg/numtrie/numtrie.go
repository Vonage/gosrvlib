/*
Package numtrie provides a numerical-indexed trie data structure. The trie
allows to store and retrieve values of any type associated with a numerical key.
*/
package numtrie

import (
	"github.com/Vonage/gosrvlib/pkg/phonekeypad"
)

const indexSize = 10 // digits from 0 to 9

// Status codes to be returned when searching for a number in the trie.
const (
	// MatchStatusNo indicates that no match was found. The first number digit
	// doesn't match any value at the trie root.
	MatchStatusNo int8 = -1

	// MatchStatusOK indicates that a full exact match was found. The full
	// number matches a trie leaf.
	MatchStatusOK int8 = 0

	// MatchStatusPrefix indicates that only a prefix of the number matches a
	// trie leaf. The remaining digits are not present in the trie.
	MatchStatusPrefix int8 = 1

	// MatchStatusPartial indicates that the full number matches a trie node
	// that is not a leaf.
	MatchStatusPartial int8 = 2

	// MatchStatusPartialPrefix indicates that only a prefix of the number
	// matches a trie node that is not a leaf. The remaining digits are not
	// present in the trie.
	MatchStatusPartialPrefix int8 = 4
)

// Node is a numerical-indexed trie node that stores a value of any type.
type Node[T any] struct {
	value       *T
	numChildren int
	children    [indexSize]*Node[T]
}

// New creates a new Node.
func New[T any]() *Node[T] {
	return &Node[T]{}
}

// Add adds a value to the trie with the given numerical key.
// It returns false if the value overrides an existing non-nil one.
func (t *Node[T]) Add(num string, val *T) bool {
	node := t

	for _, v := range num {
		i, ok := phonekeypad.KeypadDigit(v)
		if !ok {
			continue
		}

		if node.children[i] == nil {
			node.children[i] = New[T]()
			node.numChildren++
		}

		node = node.children[i]
	}

	isnew := (node.value == nil)

	node.value = val

	return isnew
}

// Get retrieves a value from the trie with the given numerical key.
// It supports partial matches.
// The return value should always be checked for the nil value.
// The second return value provides information about the match status:
//   - MatchStatusNo indicates that no match was found. The first number digit
//     doesn't match any value at the trie root.
//   - MatchStatusOK indicates that a full exact match was found. The full number
//     matches a trie leaf.
//   - MatchStatusPrefix indicates that only a prefix of the number matches a
//     trie leaf. The remaining digits are not present in the trie.
//   - MatchStatusPartial indicates that the full number matches a trie node that
//     is not a leaf.
func (t *Node[T]) Get(num string) (*T, int8) {
	node := t

	var match, digit int

	for _, v := range num {
		i, ok := phonekeypad.KeypadDigit(v)
		if !ok {
			continue
		}

		digit++

		if node.children[i] == nil {
			break
		}

		node = node.children[i]

		match++
	}

	if match == 0 {
		return node.value, MatchStatusNo
	}

	isLeaf := (node.numChildren == 0)

	if digit == match {
		if isLeaf {
			return node.value, MatchStatusOK
		}

		return node.value, MatchStatusPartial
	}

	if isLeaf {
		return node.value, MatchStatusPrefix
	}

	return node.value, MatchStatusPartialPrefix
}
