/*
Package countryphone provides geographical information of phone numbers.
Country phone codes are defined by the International Telecommunication Union
(ITU) in ITU-T standards E.123 and E.164.
*/
package countryphone

import (
	"fmt"

	"github.com/Vonage/gosrvlib/pkg/numtrie"
)

// GeoInfo stores geographical information of a phone number.
type GeoInfo struct {
	// Alpha2 is the ISO-3166 Alpha-2 Country Code.
	Alpha2 string `json:"alpha2"`

	// Area is the geographical area.
	Area string `json:"area"`

	// Type is the type of area:
	//   - 0 = ""
	//   - 1 = "state"
	//   - 2 = "province or territory"
	//   - 3 = "nation or territory"
	//   - 4 = "non-geographic"
	//   - 5 = "other"
	Type int `json:"type"`
}

// NumInfo stores the number type and geographical information of a phone number.
type NumInfo struct {
	// Type is the type of number:
	//   - 0 = ""
	//   - 1 = "landline"
	//   - 2 = "mobile"
	//   - 3 = "pager"
	//   - 4 = "satellite"
	//   - 5 = "special service"
	//   - 6 = "virtual"
	//   - 7 = "other"
	Type int `json:"type"`

	// Geo is the geographical information.
	Geo []*GeoInfo `json:"geo"`
}

// PrefixData is a type alias for a map of phone number prefixes to NumData.
type PrefixData = map[string]*NumInfo

// Data is the main data structure that stores phone number prefixes and their
// information.
type Data struct {
	enumNumberType [8]string
	enumAreaType   [6]string
	trie           *numtrie.Node[NumInfo]
}

// New initialize the search trie with the given data.
// If data is nil, the embedded default dataset is used.
func New(data PrefixData) *Data {
	d := &Data{}

	d.loadEnums()

	if data == nil {
		data = defaultData()
	}

	d.loadData(data)

	return d
}

// NumberInfo returns the number type and geographical information for the given
// phone number prefix.
//
// NOTE: see the "github.com/Vonage/gosrvlib/pkg/countrycode" package to get the
// country information from the Alpha2 code.
func (d *Data) NumberInfo(num string) (*NumInfo, error) {
	data, status := d.trie.Get(num)

	if status < 0 || data == nil {
		return nil, fmt.Errorf("no match for prefix %s", num)
	}

	return data, nil
}

// NumberType returns the string representation of the number type.
func (d *Data) NumberType(t int) (string, error) {
	if t < 0 || t >= len(d.enumNumberType) {
		return "", fmt.Errorf("invalid number type %d", t)
	}

	return d.enumNumberType[t], nil
}

// AreaType returns the string representation of the area type.
func (d *Data) AreaType(t int) (string, error) {
	if t < 0 || t >= len(d.enumAreaType) {
		return "", fmt.Errorf("invalid area type %d", t)
	}

	return d.enumAreaType[t], nil
}

// loadEnums initializes the enumeration arrays.
func (d *Data) loadEnums() {
	d.enumNumberType = [...]string{
		"",
		"landline",
		"mobile",
		"pager",
		"satellite",
		"special service",
		"virtual",
		"other",
	}

	d.enumAreaType = [...]string{
		"",
		"state",
		"province or territory",
		"nation or territory",
		"non-geographic",
		"other",
	}
}

// loadData loads the phone number prefixes and their data into the trie.
func (d *Data) loadData(data PrefixData) {
	d.trie = numtrie.New[NumInfo]()

	for k, v := range data {
		d.trie.Add(k, v)
	}
}

// defaultData returns a default map of phone numbers to country ISO-3166 Alpha-2 Codes.
// Ref.: https://en.wikipedia.org/wiki/List_of_country_calling_codes
func defaultData() PrefixData {
	return PrefixData{}
}
