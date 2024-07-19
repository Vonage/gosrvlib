package countrycode_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Vonage/gosrvlib/pkg/countrycode"
)

func ExampleData_CountryByAlpha2Code() {
	data := countrycode.New()

	got, err := data.CountryByAlpha2Code("ZW")
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	// Output:
	// {
	//   "statusCode": 1,
	//   "status": "Officially assigned",
	//   "alpha2Code": "ZW",
	//   "alpha3Code": "ZWE",
	//   "numericCode": "716",
	//   "nameEnglish": "Zimbabwe",
	//   "nameFrench": "Zimbabwe (le)",
	//   "region": "Africa",
	//   "subRegion": "Sub-Saharan Africa",
	//   "intermediateRegion": "Eastern Africa",
	//   "regionCode": "002",
	//   "subRegionCode": "202",
	//   "intermediateRegionCode": "014",
	//   "tld": "zw"
	// }
}

func ExampleData_CountryByAlpha3Code() {
	data := countrycode.New()

	got, err := data.CountryByAlpha3Code("ZWE")
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	// Output:
	// {
	//   "statusCode": 1,
	//   "status": "Officially assigned",
	//   "alpha2Code": "ZW",
	//   "alpha3Code": "ZWE",
	//   "numericCode": "716",
	//   "nameEnglish": "Zimbabwe",
	//   "nameFrench": "Zimbabwe (le)",
	//   "region": "Africa",
	//   "subRegion": "Sub-Saharan Africa",
	//   "intermediateRegion": "Eastern Africa",
	//   "regionCode": "002",
	//   "subRegionCode": "202",
	//   "intermediateRegionCode": "014",
	//   "tld": "zw"
	// }
}

func ExampleData_CountryByNumericCode() {
	data := countrycode.New()

	got, err := data.CountryByNumericCode("716")
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	// Output:
	// {
	//   "statusCode": 1,
	//   "status": "Officially assigned",
	//   "alpha2Code": "ZW",
	//   "alpha3Code": "ZWE",
	//   "numericCode": "716",
	//   "nameEnglish": "Zimbabwe",
	//   "nameFrench": "Zimbabwe (le)",
	//   "region": "Africa",
	//   "subRegion": "Sub-Saharan Africa",
	//   "intermediateRegion": "Eastern Africa",
	//   "regionCode": "002",
	//   "subRegionCode": "202",
	//   "intermediateRegionCode": "014",
	//   "tld": "zw"
	// }
}
