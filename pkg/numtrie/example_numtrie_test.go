package numtrie_test

import (
	"fmt"

	"github.com/Vonage/gosrvlib/pkg/numtrie"
)

func ExampleNode_Get() {
	// create a new numerical-indexed trie that holds sting values
	node := numtrie.New[string]()

	valA := "gamma"
	node.Add("702", &valA)

	valB := "foxtrot"
	node.Add("702153", &valB)

	// MatchStatusNo (-1) indicates that no match was found. The first number digit
	// doesn't match any value at the trie root.
	got, status := node.Get("111")
	if got != nil {
		fmt.Println(*got, status)
	} else {
		fmt.Println(got, status)
	}

	// MatchStatusOK (0) indicates that a full exact match was found. The full
	// number matches a trie leaf.
	got, status = node.Get("702153")
	if got != nil {
		fmt.Println(*got, status)
	}

	// MatchStatusPrefix (1) indicates that only a prefix of the number matches a
	// trie leaf. The remaining digits are not present in the trie.
	got, status = node.Get("702153-99")
	if got != nil {
		fmt.Println(*got, status)
	}

	// MatchStatusPartial (2) indicates that the full number matches a trie node
	// that is not a leaf.
	got, status = node.Get("702")
	if got != nil {
		fmt.Println(*got, status)
	}

	// MatchStatusPartialPrefix (4) indicates that only a prefix of the number
	// matches a trie node that is not a leaf. The remaining digits are not
	// present in the trie.
	got, status = node.Get("702-99")
	if got != nil {
		fmt.Println(*got, status)
	}

	// Output:
	// <nil> -1
	// foxtrot 0
	// foxtrot 1
	// gamma 2
	// gamma 4
}
