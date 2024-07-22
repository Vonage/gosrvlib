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

//nolint:testableexamples
func ExampleData_CountryByAlpha2Code_export() {
	data := countrycode.New()

	expCountries := make([]*countrycode.CountryData, 0, 26*26)

	// Generate all 2-letter country codes possible combinations, from AA to ZZ.
	for c1 := 'A'; c1 <= 'Z'; c1++ {
		for c0 := 'A'; c0 <= 'Z'; c0++ {
			alpha2 := string([]rune{c1, c0})

			// Decode country data from the 2-letter country code.
			country, err := data.CountryByAlpha2Code(alpha2)
			if err != nil {
				log.Fatal(err)
			}

			expCountries = append(expCountries, country)
		}
	}

	// Export the country data to JSON.
	jsonExpCountries, err := json.MarshalIndent(expCountries, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonExpCountries))
}

func ExampleData_CountryKey_import() {
	data := countrycode.New()

	// Data to import (encode).
	expCountries := []*countrycode.CountryData{
		{
			Status:                 "Officially assigned",
			Alpha2Code:             "ZM",
			Alpha3Code:             "ZMB",
			NumericCode:            "894",
			NameEnglish:            "Zambia",
			NameFrench:             "Zambie (la)",
			Region:                 "Africa",
			SubRegion:              "Sub-Saharan Africa",
			IntermediateRegion:     "Eastern Africa",
			RegionCode:             "002",
			SubRegionCode:          "202",
			IntermediateRegionCode: "014",
			TLD:                    "zm",
		},
		{
			Status:                 "Officially assigned",
			Alpha2Code:             "ZW",
			Alpha3Code:             "ZWE",
			NumericCode:            "716",
			NameEnglish:            "Zimbabwe",
			NameFrench:             "Zimbabwe (le)",
			Region:                 "Africa",
			SubRegion:              "Sub-Saharan Africa",
			IntermediateRegion:     "Eastern Africa",
			RegionCode:             "002",
			SubRegionCode:          "202",
			IntermediateRegionCode: "014",
			TLD:                    "zw",
		},
	}

	impCountryNamesByAlpha2ID := make(map[uint16]*countrycode.Names, len(expCountries))
	impCountryKeyByAlpha2ID := make(map[uint16]uint64, len(expCountries))

	for _, country := range expCountries {
		a2, countryKey, err := data.CountryKey(country)
		if err != nil {
			log.Fatal(err)
		}

		impCountryKeyByAlpha2ID[a2] = countryKey

		if len(country.NameEnglish) > 0 {
			impCountryNamesByAlpha2ID[a2] = &countrycode.Names{
				EN: country.NameEnglish,
				FR: country.NameFrench,
			}
		}
	}

	// Export the country binary data in JSON format:

	jsonCountryNamesByAlpha2ID, err := json.MarshalIndent(impCountryNamesByAlpha2ID, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonCountryNamesByAlpha2ID))

	jsonCountryKeyByAlpha2ID, err := json.MarshalIndent(impCountryKeyByAlpha2ID, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonCountryKeyByAlpha2ID))

	// Output:
	// {
	//   "845": {
	//     "EN": "Zambia",
	//     "FR": "Zambie (la)"
	//   },
	//   "855": {
	//     "EN": "Zimbabwe",
	//     "FR": "Zimbabwe (le)"
	//   }
	// }
	// {
	//   "845": 2105236111933051725,
	//   "855": 2116506203224281943
	// }
}
