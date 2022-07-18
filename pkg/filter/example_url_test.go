package filter

import (
	"fmt"
	"log"
	"net/url"
)

type Address struct {
	Country string `json:"country"`
}

type ID struct {
	Name string  `json:"name"`
	Age  int     `json:"age"`
	Addr Address `json:"address"`
}

func ExampleNew_fromURL() {
	// Simulate an encoded query passed in the http.Request of a http.Handler
	encodedJSONFilter := "%5B%5B%7B%22field%22%3A%22name%22%2C%22type%22%3A%22exact%22%2C%22value%22%3A%22doe%22%7D%2C%7B%22field%22%3A%22age%22%2C%22type%22%3A%22exact%22%2C%22value%22%3A42%7D%5D%2C%5B%7B%22field%22%3A%22address.country%22%2C%22type%22%3A%22regexp%22%2C%22value%22%3A%22UK%7CFR%22%7D%5D%5D"

	u, err := url.Parse("https://server.com/items?filter=" + encodedJSONFilter)
	if err != nil {
		log.Fatal(err)
	}

	// The filter matches the following pretty printed json:
	// [
	//   [
	//     {
	//       "field": "name",
	//       "type": "exact",
	//       "value": "doe"
	//     },
	//     {
	//       "field": "age",
	//       "type": "exact",
	//       "value": 42
	//     }
	//   ],
	//   [
	//     {
	//       "field": "address.country",
	//       "type": "regexp",
	//       "value": "EN|FR"
	//     }
	//   ]
	// ]
	// It means that either the name OR the age must match exactly AND the country must match its regular expression.
	rawFilter := GetFilter(u)
	fmt.Println(rawFilter)

	rules, err := ParseRules(rawFilter)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the filter with options
	// * WithJSONValues: We want to be lenient on the typing since we create the filter from JSON which handles a few types
	// * WithFieldNameTag: to express the filter based on JSON tags and not the actual field names
	f, err := New(
		rules,
		WithFieldNameTag("json"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Given this list, the last item will be filtered
	list := []ID{
		{
			Name: "doe",
			Age:  35,
			Addr: Address{
				Country: "UK",
			},
		},
		{
			Name: "dupont",
			Age:  42,
			Addr: Address{
				Country: "FR",
			},
		},
		{
			Name: "doe",
			Age:  42,
			Addr: Address{
				Country: "US",
			},
		},
	}

	// Filters the list in place
	err = f.Apply(&list)
	if err != nil {
		log.Fatal(err)
	}

	for _, id := range list {
		fmt.Println(id)
	}

	// Output:
	// [[{"field":"name","type":"exact","value":"doe"},{"field":"age","type":"exact","value":42}],[{"field":"address.country","type":"regexp","value":"UK|FR"}]]
	// {doe 35 {UK}}
	// {dupont 42 {FR}}
}
