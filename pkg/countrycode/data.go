package countrycode

import "strings"

// Constants for the enumeration indexes.
const (
	lenEnumStatus    = 7
	lenEnumRegion    = 6
	lenEnumSubRegion = 18
	lenEnumIntRegion = 8
)

// enumData contains the code and name of an enumeration.
type enumData struct {
	code string
	name string
}

// names contains the English and French names of a country.
type names struct {
	en string
	fr string
}

// Data contains the internal country Data and various indexes.
type Data struct {
	dStatusByID                      [lenEnumStatus]*enumData
	dStatusIDByName                  map[string]uint8
	dRegionByID                      [lenEnumRegion]*enumData
	dRegionIDByCode                  map[string]uint8
	dRegionIDByName                  map[string]uint8
	dSubRegionByID                   [lenEnumSubRegion]*enumData
	dSubRegionIDByCode               map[string]uint8
	dSubRegionIDByName               map[string]uint8
	dIntermediateRegionByID          [lenEnumIntRegion]*enumData
	dIntermediateRegionIDByCode      map[string]uint8
	dIntermediateRegionIDByName      map[string]uint8
	dCountryNamesByAlpha2ID          map[uint16]*names
	dCountryKeyByAlpha2ID            map[uint16]uint64
	dAlpha2IDByAlpha3ID              map[uint16]uint16
	dAlpha2IDByNumericID             map[uint16]uint16
	dAlpha2IDsByRegionID             map[uint8][]uint16
	dAlpha2IDsBySubRegionID          map[uint8][]uint16
	dAlpha2IDsByIntermediateRegionID map[uint8][]uint16
	dAlpha2IDsByStatusID             map[uint8][]uint16
	dAlpha2IDsByTLD                  map[uint16][]uint16
}

// New generates the country data, including various indexes.
// The generated object should be reused to avoid regenerating the data.
//
// Data sources:
//   - https://www.iso.org/iso-3166-country-codes.html
//   - https://www.cia.gov/the-world-factbook/references/country-data-codes/
//   - https://unstats.un.org/unsd/methodology/m49/overview/
//   - https://en.wikipedia.org/wiki/ISO_3166
//   - https://en.wikipedia.org/wiki/ISO_3166-1
//   - https://en.wikipedia.org/wiki/ISO_3166-2
//   - https://en.wikipedia.org/wiki/ISO_3166-3
//   - https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
//   - https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3
//   - https://en.wikipedia.org/wiki/ISO_3166-1_numeric
//
// Data updated at: 2024-07-17.
//
//nolint:funlen,maintidx
func New() *Data {
	d := &Data{
		dStatusByID: [...]*enumData{
			{"0", "Unassigned"},
			{"1", "Officially assigned"},
			{"2", "User-assigned"},
			{"3", "Exceptionally reserved"},
			{"4", "Transitionally reserved"},
			{"5", "Indeterminately reserved"},
			{"6", "Formerly assigned"},
		},
		dRegionByID: [...]*enumData{
			{"", ""},
			{"002", "Africa"},
			{"009", "Oceania"},
			{"019", "Americas"},
			{"142", "Asia"},
			{"150", "Europe"},
		},
		dSubRegionByID: [...]*enumData{
			{"", ""},
			{"015", "Northern Africa"},
			{"021", "Northern America"},
			{"030", "Eastern Asia"},
			{"034", "Southern Asia"},
			{"035", "South-eastern Asia"},
			{"039", "Southern Europe"},
			{"053", "Australia and New Zealand"},
			{"054", "Melanesia"},
			{"057", "Micronesia"},
			{"061", "Polynesia"},
			{"143", "Central Asia"},
			{"145", "Western Asia"},
			{"151", "Eastern Europe"},
			{"154", "Northern Europe"},
			{"155", "Western Europe"},
			{"202", "Sub-Saharan Africa"},
			{"419", "Latin America and the Caribbean"},
		},
		dIntermediateRegionByID: [...]*enumData{
			{"", ""},
			{"005", "South America"},
			{"011", "Western Africa"},
			{"013", "Central America"},
			{"014", "Eastern Africa"},
			{"017", "Middle Africa"},
			{"018", "Southern Africa"},
			{"029", "Caribbean"},
		},
		// The key is the alpha-2 ID encoded in 10 bits (see encodeAlpha2).
		dCountryNamesByAlpha2ID: map[uint16]*names{
			0x0024: {
				en: "Andorra",
				fr: "Andorre (l')",
			},
			0x0025: {
				en: "United Arab Emirates (the)",
				fr: "Émirats arabes unis (les)",
			},
			0x0026: {
				en: "Afghanistan",
				fr: "Afghanistan (l')",
			},
			0x0027: {
				en: "Antigua and Barbuda",
				fr: "Antigua-et-Barbuda",
			},
			0x0029: {
				en: "Anguilla",
				fr: "Anguilla",
			},
			0x002C: {
				en: "Albania",
				fr: "Albanie (l')",
			},
			0x002D: {
				en: "Armenia",
				fr: "Arménie (l')",
			},
			0x002F: {
				en: "Angola",
				fr: "Angola (l')",
			},
			0x0031: {
				en: "Antarctica",
				fr: "Antarctique (l')",
			},
			0x0032: {
				en: "Argentina",
				fr: "Argentine (l')",
			},
			0x0033: {
				en: "American Samoa",
				fr: "Samoa américaines (les)",
			},
			0x0034: {
				en: "Austria",
				fr: "Autriche (l')",
			},
			0x0035: {
				en: "Australia",
				fr: "Australie (l')",
			},
			0x0037: {
				en: "Aruba",
				fr: "Aruba",
			},
			0x0038: {
				en: "Åland Islands",
				fr: "Åland(les Îles)",
			},
			0x003A: {
				en: "Azerbaijan",
				fr: "Azerbaïdjan (l')",
			},
			0x0041: {
				en: "Bosnia and Herzegovina",
				fr: "Bosnie-Herzégovine (la)",
			},
			0x0042: {
				en: "Barbados",
				fr: "Barbade (la)",
			},
			0x0044: {
				en: "Bangladesh",
				fr: "Bangladesh (le)",
			},
			0x0045: {
				en: "Belgium",
				fr: "Belgique (la)",
			},
			0x0046: {
				en: "Burkina Faso",
				fr: "Burkina Faso (le)",
			},
			0x0047: {
				en: "Bulgaria",
				fr: "Bulgarie (la)",
			},
			0x0048: {
				en: "Bahrain",
				fr: "Bahreïn",
			},
			0x0049: {
				en: "Burundi",
				fr: "Burundi (le)",
			},
			0x004A: {
				en: "Benin",
				fr: "Bénin (le)",
			},
			0x004C: {
				en: "Saint Barthélemy",
				fr: "Saint-Barthélemy",
			},
			0x004D: {
				en: "Bermuda",
				fr: "Bermudes (les)",
			},
			0x004E: {
				en: "Brunei Darussalam",
				fr: "Brunéi Darussalam (le)",
			},
			0x004F: {
				en: "Bolivia (Plurinational State of)",
				fr: "Bolivie (État plurinational de)",
			},
			0x0051: {
				en: "Bonaire, Sint Eustatius and Saba",
				fr: "Bonaire, Saint-Eustache et Saba",
			},
			0x0052: {
				en: "Brazil",
				fr: "Brésil (le)",
			},
			0x0053: {
				en: "Bahamas (the)",
				fr: "Bahamas (les)",
			},
			0x0054: {
				en: "Bhutan",
				fr: "Bhoutan (le)",
			},
			0x0056: {
				en: "Bouvet Island",
				fr: "Bouvet (l'Île)",
			},
			0x0057: {
				en: "Botswana",
				fr: "Botswana (le)",
			},
			0x0059: {
				en: "Belarus",
				fr: "Bélarus (le)",
			},
			0x005A: {
				en: "Belize",
				fr: "Belize (le)",
			},
			0x0061: {
				en: "Canada",
				fr: "Canada (le)",
			},
			0x0063: {
				en: "Cocos (Keeling) Islands (the)",
				fr: "Cocos (les Îles)/ Keeling (les Îles)",
			},
			0x0064: {
				en: "Congo (the Democratic Republic of the)",
				fr: "Congo (la République démocratique du)",
			},
			0x0066: {
				en: "Central African Republic (the)",
				fr: "République centrafricaine (la)",
			},
			0x0067: {
				en: "Congo (the)",
				fr: "Congo (le)",
			},
			0x0068: {
				en: "Switzerland",
				fr: "Suisse (la)",
			},
			0x0069: {
				en: "Côte d'Ivoire",
				fr: "Côte d'Ivoire (la)",
			},
			0x006B: {
				en: "Cook Islands (the)",
				fr: "Cook (les Îles)",
			},
			0x006C: {
				en: "Chile",
				fr: "Chili (le)",
			},
			0x006D: {
				en: "Cameroon",
				fr: "Cameroun (le)",
			},
			0x006E: {
				en: "China",
				fr: "Chine (la)",
			},
			0x006F: {
				en: "Colombia",
				fr: "Colombie (la)",
			},
			0x0072: {
				en: "Costa Rica",
				fr: "Costa Rica (le)",
			},
			0x0075: {
				en: "Cuba",
				fr: "Cuba",
			},
			0x0076: {
				en: "Cabo Verde",
				fr: "Cabo Verde",
			},
			0x0077: {
				en: "Curaçao",
				fr: "Curaçao",
			},
			0x0078: {
				en: "Christmas Island",
				fr: "Christmas (l'Île)",
			},
			0x0079: {
				en: "Cyprus",
				fr: "Chypre",
			},
			0x007A: {
				en: "Czechia",
				fr: "Tchéquie (la)",
			},
			0x0085: {
				en: "Germany",
				fr: "Allemagne (l')",
			},
			0x008A: {
				en: "Djibouti",
				fr: "Djibouti",
			},
			0x008B: {
				en: "Denmark",
				fr: "Danemark (le)",
			},
			0x008D: {
				en: "Dominica",
				fr: "Dominique (la)",
			},
			0x008F: {
				en: "Dominican Republic (the)",
				fr: "dominicaine (la République)",
			},
			0x009A: {
				en: "Algeria",
				fr: "Algérie (l')",
			},
			0x00A3: {
				en: "Ecuador",
				fr: "Équateur (l')",
			},
			0x00A5: {
				en: "Estonia",
				fr: "Estonie (l')",
			},
			0x00A7: {
				en: "Egypt",
				fr: "Égypte (l')",
			},
			0x00A8: {
				en: "Western Sahara*",
				fr: "Sahara occidental (le)*",
			},
			0x00B2: {
				en: "Eritrea",
				fr: "Érythrée (l')",
			},
			0x00B3: {
				en: "Spain",
				fr: "Espagne (l')",
			},
			0x00B4: {
				en: "Ethiopia",
				fr: "Éthiopie (l')",
			},
			0x00C9: {
				en: "Finland",
				fr: "Finlande (la)",
			},
			0x00CA: {
				en: "Fiji",
				fr: "Fidji (les)",
			},
			0x00CB: {
				en: "Falkland Islands (the) [Malvinas]",
				fr: "Falkland (les Îles)/Malouines (les Îles)",
			},
			0x00CD: {
				en: "Micronesia (Federated States of)",
				fr: "Micronésie (États fédérés de)",
			},
			0x00CF: {
				en: "Faroe Islands (the)",
				fr: "Féroé (les Îles)",
			},
			0x00D2: {
				en: "France",
				fr: "France (la)",
			},
			0x00E1: {
				en: "Gabon",
				fr: "Gabon (le)",
			},
			0x00E2: {
				en: "United Kingdom of Great Britain and Northern Ireland (the)",
				fr: "Royaume-Uni de Grande-Bretagne et d'Irlande du Nord (le)",
			},
			0x00E4: {
				en: "Grenada",
				fr: "Grenade (la)",
			},
			0x00E5: {
				en: "Georgia",
				fr: "Géorgie (la)",
			},
			0x00E6: {
				en: "French Guiana",
				fr: "Guyane française (la )",
			},
			0x00E7: {
				en: "Guernsey",
				fr: "Guernesey",
			},
			0x00E8: {
				en: "Ghana",
				fr: "Ghana (le)",
			},
			0x00E9: {
				en: "Gibraltar",
				fr: "Gibraltar",
			},
			0x00EC: {
				en: "Greenland",
				fr: "Groenland (le)",
			},
			0x00ED: {
				en: "Gambia (the)",
				fr: "Gambie (la)",
			},
			0x00EE: {
				en: "Guinea",
				fr: "Guinée (la)",
			},
			0x00F0: {
				en: "Guadeloupe",
				fr: "Guadeloupe (la)",
			},
			0x00F1: {
				en: "Equatorial Guinea",
				fr: "Guinée équatoriale (la)",
			},
			0x00F2: {
				en: "Greece",
				fr: "Grèce (la)",
			},
			0x00F3: {
				en: "South Georgia and the South Sandwich Islands",
				fr: "Géorgie du Sud-et-les Îles Sandwich du Sud (la)",
			},
			0x00F4: {
				en: "Guatemala",
				fr: "Guatemala (le)",
			},
			0x00F5: {
				en: "Guam",
				fr: "Guam",
			},
			0x00F7: {
				en: "Guinea-Bissau",
				fr: "Guinée-Bissau (la)",
			},
			0x00F9: {
				en: "Guyana",
				fr: "Guyana (le)",
			},
			0x010B: {
				en: "Hong Kong",
				fr: "Hong Kong",
			},
			0x010D: {
				en: "Heard Island and McDonald Islands",
				fr: "Heard-et-Îles MacDonald (l'Île)",
			},
			0x010E: {
				en: "Honduras",
				fr: "Honduras (le)",
			},
			0x0112: {
				en: "Croatia",
				fr: "Croatie (la)",
			},
			0x0114: {
				en: "Haiti",
				fr: "Haïti",
			},
			0x0115: {
				en: "Hungary",
				fr: "Hongrie (la)",
			},
			0x0124: {
				en: "Indonesia",
				fr: "Indonésie (l')",
			},
			0x0125: {
				en: "Ireland",
				fr: "Irlande (l')",
			},
			0x012C: {
				en: "Israel",
				fr: "Israël",
			},
			0x012D: {
				en: "Isle of Man",
				fr: "Île de Man",
			},
			0x012E: {
				en: "India",
				fr: "Inde (l')",
			},
			0x012F: {
				en: "British Indian Ocean Territory (the)",
				fr: "Indien (le Territoire britannique de l'océan)",
			},
			0x0131: {
				en: "Iraq",
				fr: "Iraq (l')",
			},
			0x0132: {
				en: "Iran (Islamic Republic of)",
				fr: "Iran (République Islamique d')",
			},
			0x0133: {
				en: "Iceland",
				fr: "Islande (l')",
			},
			0x0134: {
				en: "Italy",
				fr: "Italie (l')",
			},
			0x0145: {
				en: "Jersey",
				fr: "Jersey",
			},
			0x014D: {
				en: "Jamaica",
				fr: "Jamaïque (la)",
			},
			0x014F: {
				en: "Jordan",
				fr: "Jordanie (la)",
			},
			0x0150: {
				en: "Japan",
				fr: "Japon (le)",
			},
			0x0165: {
				en: "Kenya",
				fr: "Kenya (le)",
			},
			0x0167: {
				en: "Kyrgyzstan",
				fr: "Kirghizistan (le)",
			},
			0x0168: {
				en: "Cambodia",
				fr: "Cambodge (le)",
			},
			0x0169: {
				en: "Kiribati",
				fr: "Kiribati",
			},
			0x016D: {
				en: "Comoros (the)",
				fr: "Comores (les)",
			},
			0x016E: {
				en: "Saint Kitts and Nevis",
				fr: "Saint-Kitts-et-Nevis",
			},
			0x0170: {
				en: "Korea (the Democratic People's Republic of)",
				fr: "Corée (la République populaire démocratique de)",
			},
			0x0172: {
				en: "Korea (the Republic of)",
				fr: "Corée (la République de)",
			},
			0x0177: {
				en: "Kuwait",
				fr: "Koweït (le)",
			},
			0x0179: {
				en: "Cayman Islands (the)",
				fr: "Caïmans (les Îles)",
			},
			0x017A: {
				en: "Kazakhstan",
				fr: "Kazakhstan (le)",
			},
			0x0181: {
				en: "Lao People's Democratic Republic (the)",
				fr: "Lao (la République démocratique populaire)",
			},
			0x0182: {
				en: "Lebanon",
				fr: "Liban (le)",
			},
			0x0183: {
				en: "Saint Lucia",
				fr: "Sainte-Lucie",
			},
			0x0189: {
				en: "Liechtenstein",
				fr: "Liechtenstein (le)",
			},
			0x018B: {
				en: "Sri Lanka",
				fr: "Sri Lanka",
			},
			0x0192: {
				en: "Liberia",
				fr: "Libéria (le)",
			},
			0x0193: {
				en: "Lesotho",
				fr: "Lesotho (le)",
			},
			0x0194: {
				en: "Lithuania",
				fr: "Lituanie (la)",
			},
			0x0195: {
				en: "Luxembourg",
				fr: "Luxembourg (le)",
			},
			0x0196: {
				en: "Latvia",
				fr: "Lettonie (la)",
			},
			0x0199: {
				en: "Libya",
				fr: "Libye (la)",
			},
			0x01A1: {
				en: "Morocco",
				fr: "Maroc (le)",
			},
			0x01A3: {
				en: "Monaco",
				fr: "Monaco",
			},
			0x01A4: {
				en: "Moldova (the Republic of)",
				fr: "Moldova (la République de)",
			},
			0x01A5: {
				en: "Montenegro",
				fr: "Monténégro (le)",
			},
			0x01A6: {
				en: "Saint Martin (French part)",
				fr: "Saint-Martin (partie française)",
			},
			0x01A7: {
				en: "Madagascar",
				fr: "Madagascar",
			},
			0x01A8: {
				en: "Marshall Islands (the)",
				fr: "Marshall (les Îles)",
			},
			0x01AB: {
				en: "North Macedonia",
				fr: "Macédoine du Nord (la)",
			},
			0x01AC: {
				en: "Mali",
				fr: "Mali (le)",
			},
			0x01AD: {
				en: "Myanmar",
				fr: "Myanmar (le)",
			},
			0x01AE: {
				en: "Mongolia",
				fr: "Mongolie (la)",
			},
			0x01AF: {
				en: "Macao",
				fr: "Macao",
			},
			0x01B0: {
				en: "Northern Mariana Islands (the)",
				fr: "Mariannes du Nord (les Îles)",
			},
			0x01B1: {
				en: "Martinique",
				fr: "Martinique (la)",
			},
			0x01B2: {
				en: "Mauritania",
				fr: "Mauritanie (la)",
			},
			0x01B3: {
				en: "Montserrat",
				fr: "Montserrat",
			},
			0x01B4: {
				en: "Malta",
				fr: "Malte",
			},
			0x01B5: {
				en: "Mauritius",
				fr: "Maurice",
			},
			0x01B6: {
				en: "Maldives",
				fr: "Maldives (les)",
			},
			0x01B7: {
				en: "Malawi",
				fr: "Malawi (le)",
			},
			0x01B8: {
				en: "Mexico",
				fr: "Mexique (le)",
			},
			0x01B9: {
				en: "Malaysia",
				fr: "Malaisie (la)",
			},
			0x01BA: {
				en: "Mozambique",
				fr: "Mozambique (le)",
			},
			0x01C1: {
				en: "Namibia",
				fr: "Namibie (la)",
			},
			0x01C3: {
				en: "New Caledonia",
				fr: "Nouvelle-Calédonie (la)",
			},
			0x01C5: {
				en: "Niger (the)",
				fr: "Niger (le)",
			},
			0x01C6: {
				en: "Norfolk Island",
				fr: "Norfolk (l'Île)",
			},
			0x01C7: {
				en: "Nigeria",
				fr: "Nigéria (le)",
			},
			0x01C9: {
				en: "Nicaragua",
				fr: "Nicaragua (le)",
			},
			0x01CC: {
				en: "Netherlands (Kingdom of the)",
				fr: "Pays-Bas (Royaume des)",
			},
			0x01CF: {
				en: "Norway",
				fr: "Norvège (la)",
			},
			0x01D0: {
				en: "Nepal",
				fr: "Népal (le)",
			},
			0x01D2: {
				en: "Nauru",
				fr: "Nauru",
			},
			0x01D5: {
				en: "Niue",
				fr: "Niue",
			},
			0x01DA: {
				en: "New Zealand",
				fr: "Nouvelle-Zélande (la)",
			},
			0x01ED: {
				en: "Oman",
				fr: "Oman",
			},
			0x0201: {
				en: "Panama",
				fr: "Panama (le)",
			},
			0x0205: {
				en: "Peru",
				fr: "Pérou (le)",
			},
			0x0206: {
				en: "French Polynesia",
				fr: "Polynésie française (la)",
			},
			0x0207: {
				en: "Papua New Guinea",
				fr: "Papouasie-Nouvelle-Guinée (la)",
			},
			0x0208: {
				en: "Philippines (the)",
				fr: "Philippines (les)",
			},
			0x020B: {
				en: "Pakistan",
				fr: "Pakistan (le)",
			},
			0x020C: {
				en: "Poland",
				fr: "Pologne (la)",
			},
			0x020D: {
				en: "Saint Pierre and Miquelon",
				fr: "Saint-Pierre-et-Miquelon",
			},
			0x020E: {
				en: "Pitcairn",
				fr: "Pitcairn",
			},
			0x0212: {
				en: "Puerto Rico",
				fr: "Porto Rico",
			},
			0x0213: {
				en: "Palestine, State of",
				fr: "Palestine, État de",
			},
			0x0214: {
				en: "Portugal",
				fr: "Portugal (le)",
			},
			0x0217: {
				en: "Palau",
				fr: "Palaos (les)",
			},
			0x0219: {
				en: "Paraguay",
				fr: "Paraguay (le)",
			},
			0x0221: {
				en: "Qatar",
				fr: "Qatar (le)",
			},
			0x0245: {
				en: "Réunion",
				fr: "Réunion (La)",
			},
			0x024F: {
				en: "Romania",
				fr: "Roumanie (la)",
			},
			0x0253: {
				en: "Serbia",
				fr: "Serbie (la)",
			},
			0x0255: {
				en: "Russian Federation (the)",
				fr: "Russie (la Fédération de)",
			},
			0x0257: {
				en: "Rwanda",
				fr: "Rwanda (le)",
			},
			0x0261: {
				en: "Saudi Arabia",
				fr: "Arabie saoudite (l')",
			},
			0x0262: {
				en: "Solomon Islands",
				fr: "Salomon (les Îles)",
			},
			0x0263: {
				en: "Seychelles",
				fr: "Seychelles (les)",
			},
			0x0264: {
				en: "Sudan (the)",
				fr: "Soudan (le)",
			},
			0x0265: {
				en: "Sweden",
				fr: "Suède (la)",
			},
			0x0267: {
				en: "Singapore",
				fr: "Singapour",
			},
			0x0268: {
				en: "Saint Helena, Ascension and Tristan da Cunha",
				fr: "Sainte-Hélène, Ascension et Tristan da Cunha",
			},
			0x0269: {
				en: "Slovenia",
				fr: "Slovénie (la)",
			},
			0x026A: {
				en: "Svalbard and Jan Mayen",
				fr: "Svalbard et l'Île Jan Mayen (le)",
			},
			0x026B: {
				en: "Slovakia",
				fr: "Slovaquie (la)",
			},
			0x026C: {
				en: "Sierra Leone",
				fr: "Sierra Leone (la)",
			},
			0x026D: {
				en: "San Marino",
				fr: "Saint-Marin",
			},
			0x026E: {
				en: "Senegal",
				fr: "Sénégal (le)",
			},
			0x026F: {
				en: "Somalia",
				fr: "Somalie (la)",
			},
			0x0272: {
				en: "Suriname",
				fr: "Suriname (le)",
			},
			0x0273: {
				en: "South Sudan",
				fr: "Soudan du Sud (le)",
			},
			0x0274: {
				en: "Sao Tome and Principe",
				fr: "Sao Tomé-et-Principe",
			},
			0x0276: {
				en: "El Salvador",
				fr: "El Salvador",
			},
			0x0278: {
				en: "Sint Maarten (Dutch part)",
				fr: "Saint-Martin (partie néerlandaise)",
			},
			0x0279: {
				en: "Syrian Arab Republic (the)",
				fr: "République arabe syrienne (la)",
			},
			0x027A: {
				en: "Eswatini",
				fr: "Eswatini (l')",
			},
			0x0283: {
				en: "Turks and Caicos Islands (the)",
				fr: "Turks-et-Caïcos (les Îles)",
			},
			0x0284: {
				en: "Chad",
				fr: "Tchad (le)",
			},
			0x0286: {
				en: "French Southern Territories (the)",
				fr: "Terres australes françaises (les)",
			},
			0x0287: {
				en: "Togo",
				fr: "Togo (le)",
			},
			0x0288: {
				en: "Thailand",
				fr: "Thaïlande (la)",
			},
			0x028A: {
				en: "Tajikistan",
				fr: "Tadjikistan (le)",
			},
			0x028B: {
				en: "Tokelau",
				fr: "Tokelau (les)",
			},
			0x028C: {
				en: "Timor-Leste",
				fr: "Timor-Leste (le)",
			},
			0x028D: {
				en: "Turkmenistan",
				fr: "Turkménistan (le)",
			},
			0x028E: {
				en: "Tunisia",
				fr: "Tunisie (la)",
			},
			0x028F: {
				en: "Tonga",
				fr: "Tonga (les)",
			},
			0x0292: {
				en: "Türkiye",
				fr: "Türkiye (la)",
			},
			0x0294: {
				en: "Trinidad and Tobago",
				fr: "Trinité-et-Tobago (la)",
			},
			0x0296: {
				en: "Tuvalu",
				fr: "Tuvalu (les)",
			},
			0x0297: {
				en: "Taiwan (Province of China)",
				fr: "Taïwan (Province de Chine)",
			},
			0x029A: {
				en: "Tanzania, the United Republic of",
				fr: "Tanzanie (la République-Unie de)",
			},
			0x02A1: {
				en: "Ukraine",
				fr: "Ukraine (l')",
			},
			0x02A7: {
				en: "Uganda",
				fr: "Ouganda (l')",
			},
			0x02AD: {
				en: "United States Minor Outlying Islands (the)",
				fr: "Îles mineures éloignées des États-Unis (les)",
			},
			0x02B3: {
				en: "United States of America (the)",
				fr: "États-Unis d'Amérique (les)",
			},
			0x02B9: {
				en: "Uruguay",
				fr: "Uruguay (l')",
			},
			0x02BA: {
				en: "Uzbekistan",
				fr: "Ouzbékistan (l')",
			},
			0x02C1: {
				en: "Holy See (the)",
				fr: "Saint-Siège (le)",
			},
			0x02C3: {
				en: "Saint Vincent and the Grenadines",
				fr: "Saint-Vincent-et-les Grenadines",
			},
			0x02C5: {
				en: "Venezuela (Bolivarian Republic of)",
				fr: "Venezuela (République bolivarienne du)",
			},
			0x02C7: {
				en: "Virgin Islands (British)",
				fr: "Vierges britanniques (les Îles)",
			},
			0x02C9: {
				en: "Virgin Islands (U.S.)",
				fr: "Vierges des États-Unis (les Îles)",
			},
			0x02CE: {
				en: "Viet Nam",
				fr: "Viet Nam (le)",
			},
			0x02D5: {
				en: "Vanuatu",
				fr: "Vanuatu (le)",
			},
			0x02E6: {
				en: "Wallis and Futuna",
				fr: "Wallis-et-Futuna ",
			},
			0x02F3: {
				en: "Samoa",
				fr: "Samoa (le)",
			},
			0x0325: {
				en: "Yemen",
				fr: "Yémen (le)",
			},
			0x0334: {
				en: "Mayotte",
				fr: "Mayotte",
			},
			0x0341: {
				en: "South Africa",
				fr: "Afrique du Sud (l')",
			},
			0x034D: {
				en: "Zambia",
				fr: "Zambie (la)",
			},
			0x0357: {
				en: "Zimbabwe",
				fr: "Zimbabwe (le)",
			},
		},
		// dCountryKeyByAlpha2ID contains all the countries data, excluding names, in a compact form.
		// The key is the alpha-2 code encoded in 10 bits (see encodeAlpha2).
		// The value is the country data encoded in 64 bits (see decodeCountryKey).
		// This map contains all the 2-characters combinations from AA to ZZ.
		dCountryKeyByAlpha2ID: map[uint16]uint64{
			0x0021: 0x2084000000000000,
			0x0022: 0x0088000000000000,
			0x0023: 0x308C000000000000,
			0x0024: 0x10902E20294C0240,
			0x0025: 0x1094322E21180250,
			0x0026: 0x1098263809080260,
			0x0027: 0x109C343838E3C270,
			0x0028: 0x00A0000000000000,
			0x0029: 0x10A4290D28E3C290,
			0x002A: 0x00A8000000000000,
			0x002B: 0x00AC000000000000,
			0x002C: 0x10B02C10114C02C0,
			0x002D: 0x10B43268671802D0,
			0x002E: 0x40B8000000000000,
			0x002F: 0x10BC2778306142F0,
			0x0030: 0x50C0000000000000,
			0x0031: 0x10C4340814000310,
			0x0032: 0x10C8323840E24320,
			0x0033: 0x10CC336820940330,
			0x0034: 0x10D035A0515E0340,
			0x0035: 0x10D43598488E0350,
			0x0036: 0x00D8000000000000,
			0x0037: 0x10DC22BC2AE3C370,
			0x0038: 0x10E02C09F15C0380,
			0x0039: 0x00E4000000000000,
			0x003A: 0x10E83A283F1803A0,
			0x0041: 0x110449408D4C0410,
			0x0042: 0x1108521068E3C420,
			0x0043: 0x010C000000000000,
			0x0044: 0x1110472065080440,
			0x0045: 0x11144560715E0450,
			0x0046: 0x1118460EAC608460,
			0x0047: 0x111C4790C95A0470,
			0x0048: 0x1120489061180480,
			0x0049: 0x11244448D8610490,
			0x004A: 0x11284571986084A0,
			0x004B: 0x012C000000000000,
			0x004C: 0x11304C6D18E3C4C0,
			0x004D: 0x11344DA878C404D0,
			0x004E: 0x11385270C10A04E0,
			0x004F: 0x113C4F6088E244F0,
			0x0050: 0x0140000000000000,
			0x0051: 0x1144459C2EE3C510,
			0x0052: 0x1148520898E24520,
			0x0053: 0x114C489858E3C530,
			0x0054: 0x1150547081080540,
			0x0055: 0x4154000000000000,
			0x0056: 0x115856A094E24560,
			0x0057: 0x115C570890618570,
			0x0058: 0x5160000000000000,
			0x0059: 0x11644C90E15A0590,
			0x005A: 0x11684CD0A8E2C5A0,
			0x0061: 0x11846170F8C40610,
			0x0062: 0x0188000000000000,
			0x0063: 0x118C63594C8E0630,
			0x0064: 0x11906F2168614640,
			0x0065: 0x0194000000000000,
			0x0066: 0x1198613118614660,
			0x0067: 0x119C6F3964614670,
			0x0068: 0x11A0682DE95E0680,
			0x0069: 0x11A469B300608690,
			0x006A: 0x01A8000000000000,
			0x006B: 0x11AC6F59709406B0,
			0x006C: 0x11B0686130E246C0,
			0x006D: 0x11B46D90F06146D0,
			0x006E: 0x11B86871390606E0,
			0x006F: 0x11BC6F6154E246F0,
			0x0070: 0x31C0000000000000,
			0x0071: 0x31C4000000000000,
			0x0072: 0x11C8724978E2C720,
			0x0073: 0x41CC000000000000,
			0x0074: 0x61D0000000000000,
			0x0075: 0x11D4751180E3C750,
			0x0076: 0x11D870B108608760,
			0x0077: 0x11DC75BC26E3C770,
			0x0078: 0x11E07891448E0780,
			0x0079: 0x11E4798189180790,
			0x007A: 0x11E87A29975A07A0,
			0x0081: 0x0204000000000000,
			0x0082: 0x0208000000000000,
			0x0083: 0x020C000000000000,
			0x0084: 0x6210000000000000,
			0x0085: 0x121485AA295E0850,
			0x0086: 0x0218000000000000,
			0x0087: 0x321C000000000000,
			0x0088: 0x0220000000000000,
			0x0089: 0x0224000000000000,
			0x008A: 0x12288A4A0C6108A0,
			0x008B: 0x122C8E59A15C08B0,
			0x008C: 0x0230000000000000,
			0x008D: 0x12348D09A8E3C8D0,
			0x008E: 0x0238000000000000,
			0x008F: 0x123C8F69ACE3C8F0,
			0x0090: 0x0240000000000000,
			0x0091: 0x0244000000000000,
			0x0092: 0x0248000000000000,
			0x0093: 0x024C000000000000,
			0x0094: 0x0250000000000000,
			0x0095: 0x0254000000000000,
			0x0096: 0x0258000000000000,
			0x0097: 0x025C000000000000,
			0x0098: 0x0260000000000000,
			0x0099: 0x5264000000000000,
			0x009A: 0x12689A08184209A0,
			0x00A1: 0x3284000000000000,
			0x00A2: 0x0288000000000000,
			0x00A3: 0x128CA3A9B4E24A30,
			0x00A4: 0x0290000000000000,
			0x00A5: 0x1294B3A1D35C0A50,
			0x00A6: 0x5298000000000000,
			0x00A7: 0x129CA7CE64420A70,
			0x00A8: 0x12A0B345B8420A80,
			0x00A9: 0x02A4000000000000,
			0x00AA: 0x02A8000000000000,
			0x00AB: 0x02AC000000000000,
			0x00AC: 0x02B0000000000000,
			0x00AD: 0x52B4000000000000,
			0x00AE: 0x02B8000000000000,
			0x00AF: 0x02BC000000000000,
			0x00B0: 0x52C0000000000000,
			0x00B1: 0x02C4000000000000,
			0x00B2: 0x12C8B249D0610B20,
			0x00B3: 0x12CCB385A94C0B30,
			0x00B4: 0x12D0B441CE610B40,
			0x00B5: 0x32D4000000000000,
			0x00B6: 0x52D8000000000000,
			0x00B7: 0x52DC000000000000,
			0x00B8: 0x02E0000000000000,
			0x00B9: 0x02E4000000000000,
			0x00BA: 0x32E8000000000000,
			0x00C1: 0x0304000000000000,
			0x00C2: 0x0308000000000000,
			0x00C3: 0x030C000000000000,
			0x00C4: 0x0310000000000000,
			0x00C5: 0x0314000000000000,
			0x00C6: 0x0318000000000000,
			0x00C7: 0x031C000000000000,
			0x00C8: 0x0320000000000000,
			0x00C9: 0x1324C971ED5C0C90,
			0x00CA: 0x1328CA49E4900CA0,
			0x00CB: 0x132CCC59DCE24CB0,
			0x00CC: 0x5330000000000000,
			0x00CD: 0x1334D36C8E920CD0,
			0x00CE: 0x0338000000000000,
			0x00CF: 0x133CD279D55C0CF0,
			0x00D0: 0x0340000000000000,
			0x00D1: 0x6344000000000000,
			0x00D2: 0x1348D209F55E0D20,
			0x00D3: 0x034C000000000000,
			0x00D4: 0x0350000000000000,
			0x00D5: 0x0354000000000000,
			0x00D6: 0x0358000000000000,
			0x00D7: 0x035C000000000000,
			0x00D8: 0x3360000000000000,
			0x00D9: 0x0364000000000000,
			0x00DA: 0x0368000000000000,
			0x00E1: 0x1384E11214614E10,
			0x00E2: 0x1388E296755C2AB0,
			0x00E3: 0x538C000000000000,
			0x00E4: 0x1390F22268E3CE40,
			0x00E5: 0x1394E57A19180E50,
			0x00E6: 0x1398F531FCE24E60,
			0x00E7: 0x139CE7CE7F5C0E70,
			0x00E8: 0x13A0E80A40608E80,
			0x00E9: 0x13A4E912494C0E90,
			0x00EA: 0x03A8000000000000,
			0x00EB: 0x03AC000000000000,
			0x00EC: 0x13B0F26260C40EC0,
			0x00ED: 0x13B4ED121C608ED0,
			0x00EE: 0x13B8E97288608EE0,
			0x00EF: 0x03BC000000000000,
			0x00F0: 0x13C0EC8270E3CF00,
			0x00F1: 0x13C4EE89C4614F10,
			0x00F2: 0x13C8F21A594C0F20,
			0x00F3: 0x13CE6799DEE24F30,
			0x00F4: 0x13D0F46A80E2CF40,
			0x00F5: 0x13D4F56A78920F50,
			0x00F6: 0x03D8000000000000,
			0x00F7: 0x13DCEE14E0608F70,
			0x00F8: 0x03E0000000000000,
			0x00F9: 0x13E4F5CA90E24F90,
			0x00FA: 0x03E8000000000000,
			0x0101: 0x0404000000000000,
			0x0102: 0x0408000000000000,
			0x0103: 0x040C000000000000,
			0x0104: 0x0410000000000000,
			0x0105: 0x0414000000000000,
			0x0106: 0x0418000000000000,
			0x0107: 0x041C000000000000,
			0x0108: 0x0420000000000000,
			0x0109: 0x0424000000000000,
			0x010A: 0x0428000000000000,
			0x010B: 0x142D0B3AB10610B0,
			0x010C: 0x0430000000000000,
			0x010D: 0x14350D229C8E10D0,
			0x010E: 0x14390E22A8E2D0E0,
			0x010F: 0x043C000000000000,
			0x0110: 0x0440000000000000,
			0x0111: 0x0444000000000000,
			0x0112: 0x144912B17F4C1120,
			0x0113: 0x044C000000000000,
			0x0114: 0x1451144A98E3D140,
			0x0115: 0x14551572B95A1150,
			0x0116: 0x6458000000000000,
			0x0117: 0x045C000000000000,
			0x0118: 0x0460000000000000,
			0x0119: 0x0464000000000000,
			0x011A: 0x0468000000000000,
			0x0121: 0x0484000000000000,
			0x0122: 0x5488000000000000,
			0x0123: 0x348C000000000000,
			0x0124: 0x14912472D10A1240,
			0x0125: 0x14953262E95C1250,
			0x0126: 0x0498000000000000,
			0x0127: 0x049C000000000000,
			0x0128: 0x04A0000000000000,
			0x0129: 0x04A4000000000000,
			0x012A: 0x04A8000000000000,
			0x012B: 0x04AC000000000000,
			0x012C: 0x14B13392F11812C0,
			0x012D: 0x14B52D76835C12D0,
			0x012E: 0x14B92E22C90812E0,
			0x012F: 0x14BD2FA0AC6112F0,
			0x0130: 0x04C0000000000000,
			0x0131: 0x14C5328AE1181310,
			0x0132: 0x14C93272D9081320,
			0x0133: 0x14CD3362C15C1330,
			0x0134: 0x14D1340AF94C1340,
			0x0135: 0x04D4000000000000,
			0x0136: 0x04D8000000000000,
			0x0137: 0x04DC000000000000,
			0x0138: 0x04E0000000000000,
			0x0139: 0x04E4000000000000,
			0x013A: 0x04E8000000000000,
			0x0141: 0x5504000000000000,
			0x0142: 0x0508000000000000,
			0x0143: 0x050C000000000000,
			0x0144: 0x0510000000000000,
			0x0145: 0x151545CE815C1450,
			0x0146: 0x0518000000000000,
			0x0147: 0x051C000000000000,
			0x0148: 0x0520000000000000,
			0x0149: 0x0524000000000000,
			0x014A: 0x0528000000000000,
			0x014B: 0x052C000000000000,
			0x014C: 0x0530000000000000,
			0x014D: 0x1535416B08E3D4D0,
			0x014E: 0x0538000000000000,
			0x014F: 0x153D4F93211814F0,
			0x0150: 0x1541507311061500,
			0x0151: 0x0544000000000000,
			0x0152: 0x0548000000000000,
			0x0153: 0x054C000000000000,
			0x0154: 0x6550000000000000,
			0x0155: 0x0554000000000000,
			0x0156: 0x0558000000000000,
			0x0157: 0x055C000000000000,
			0x0158: 0x0560000000000000,
			0x0159: 0x0564000000000000,
			0x015A: 0x0568000000000000,
			0x0161: 0x0584000000000000,
			0x0162: 0x0588000000000000,
			0x0163: 0x058C000000000000,
			0x0164: 0x0590000000000000,
			0x0165: 0x1595657328611650,
			0x0166: 0x0598000000000000,
			0x0167: 0x159D67D343161670,
			0x0168: 0x15A16868E90A1680,
			0x0169: 0x15A5699250921690,
			0x016A: 0x05A8000000000000,
			0x016B: 0x05AC000000000000,
			0x016C: 0x05B0000000000000,
			0x016D: 0x15B46F695C6116D0,
			0x016E: 0x15B96E0D26E3D6E0,
			0x016F: 0x05BC000000000000,
			0x0170: 0x15C2125B31061700,
			0x0171: 0x05C4000000000000,
			0x0172: 0x15C96F9335061720,
			0x0173: 0x05CC000000000000,
			0x0174: 0x05D0000000000000,
			0x0175: 0x05D4000000000000,
			0x0176: 0x05D8000000000000,
			0x0177: 0x15DD77A33D181770,
			0x0178: 0x05E0000000000000,
			0x0179: 0x15E4796910E3D790,
			0x017A: 0x15E961D31D1617A0,
			0x0181: 0x1605817B450A1810,
			0x0182: 0x160982734D181820,
			0x0183: 0x160D830D2CE3D830,
			0x0184: 0x0610000000000000,
			0x0185: 0x0614000000000000,
			0x0186: 0x5618000000000000,
			0x0187: 0x061C000000000000,
			0x0188: 0x0620000000000000,
			0x0189: 0x1625892B6D5E1890,
			0x018A: 0x0628000000000000,
			0x018B: 0x162D8B09210818B0,
			0x018C: 0x0630000000000000,
			0x018D: 0x0634000000000000,
			0x018E: 0x0638000000000000,
			0x018F: 0x063C000000000000,
			0x0190: 0x0640000000000000,
			0x0191: 0x0644000000000000,
			0x0192: 0x164982935C609920,
			0x0193: 0x164D937B54619930,
			0x0194: 0x165194AB715C1940,
			0x0195: 0x165595C3755E1950,
			0x0196: 0x1659960B595C1960,
			0x0197: 0x065C000000000000,
			0x0198: 0x0660000000000000,
			0x0199: 0x166582CB64421990,
			0x019A: 0x0668000000000000,
			0x01A1: 0x1685A193F0421A10,
			0x01A2: 0x0688000000000000,
			0x01A3: 0x168DA37BD95E1A30,
			0x01A4: 0x1691A40BE55A1A40,
			0x01A5: 0x1695AE2BE74C1A50,
			0x01A6: 0x1699A1352EE3DA60,
			0x01A7: 0x169DA43B84611A70,
			0x01A8: 0x16A1A86490921A80,
			0x01A9: 0x66A4000000000000,
			0x01AA: 0x06A8000000000000,
			0x01AB: 0x16ADAB264F4C1AB0,
			0x01AC: 0x16B1AC4BA4609AC0,
			0x01AD: 0x16B5AD90D10A1AD0,
			0x01AE: 0x16B9AE3BE1061AE0,
			0x01AF: 0x16BDA11B7D061AF0,
			0x01B0: 0x16C1AE8488921B00,
			0x01B1: 0x16C5B48BB4E3DB10,
			0x01B2: 0x16C9B2A3BC609B20,
			0x01B3: 0x16CDB393E8E3DB30,
			0x01B4: 0x16D1ACA3AD4C1B40,
			0x01B5: 0x16D5B59BC0611B50,
			0x01B6: 0x16D9A4B39D081B60,
			0x01B7: 0x16DDB74B8C611B70,
			0x01B8: 0x16E1A5C3C8E2DB80,
			0x01B9: 0x16E5B99B950A1B90,
			0x01BA: 0x16E9AFD3F8611BA0,
			0x01C1: 0x1705C16C08619C10,
			0x01C2: 0x0708000000000000,
			0x01C3: 0x170DC36438901C30,
			0x01C4: 0x0710000000000000,
			0x01C5: 0x1715C59464609C50,
			0x01C6: 0x1719C65C7C8E1C60,
			0x01C7: 0x171DC70C6C609C70,
			0x01C8: 0x6720000000000000,
			0x01C9: 0x1725C91C5CE2DC90,
			0x01CA: 0x0728000000000000,
			0x01CB: 0x072C000000000000,
			0x01CC: 0x1731CC24215E1CC0,
			0x01CD: 0x0734000000000000,
			0x01CE: 0x0738000000000000,
			0x01CF: 0x173DCF94855C1CF0,
			0x01D0: 0x1741D06419081D00,
			0x01D1: 0x6744000000000000,
			0x01D2: 0x1749D2AC10921D20,
			0x01D3: 0x074C000000000000,
			0x01D4: 0x4750000000000000,
			0x01D5: 0x1755C9AC74941D50,
			0x01D6: 0x0758000000000000,
			0x01D7: 0x075C000000000000,
			0x01D8: 0x0760000000000000,
			0x01D9: 0x0764000000000000,
			0x01DA: 0x1769DA64548E1DA0,
			0x01E1: 0x5784000000000000,
			0x01E2: 0x0788000000000000,
			0x01E3: 0x078C000000000000,
			0x01E4: 0x0790000000000000,
			0x01E5: 0x0794000000000000,
			0x01E6: 0x0798000000000000,
			0x01E7: 0x079C000000000000,
			0x01E8: 0x07A0000000000000,
			0x01E9: 0x07A4000000000000,
			0x01EA: 0x07A8000000000000,
			0x01EB: 0x07AC000000000000,
			0x01EC: 0x07B0000000000000,
			0x01ED: 0x17B5ED7401181ED0,
			0x01EE: 0x07B8000000000000,
			0x01EF: 0x07BC000000000000,
			0x01F0: 0x07C0000000000000,
			0x01F1: 0x07C4000000000000,
			0x01F2: 0x07C8000000000000,
			0x01F3: 0x07CC000000000000,
			0x01F4: 0x07D0000000000000,
			0x01F5: 0x07D4000000000000,
			0x01F6: 0x07D8000000000000,
			0x01F7: 0x07DC000000000000,
			0x01F8: 0x07E0000000000000,
			0x01F9: 0x07E4000000000000,
			0x01FA: 0x07E8000000000000,
			0x0201: 0x180601749EE2E010,
			0x0202: 0x0808000000000000,
			0x0203: 0x680C000000000000,
			0x0204: 0x0810000000000000,
			0x0205: 0x18160594B8E26050,
			0x0206: 0x181A193204942060,
			0x0207: 0x181E0E3CAC902070,
			0x0208: 0x18220864C10A2080,
			0x0209: 0x5824000000000000,
			0x020A: 0x0828000000000000,
			0x020B: 0x182E015C950820B0,
			0x020C: 0x18320F64D15A20C0,
			0x020D: 0x1836706D34C420D0,
			0x020E: 0x183A0374C89420E0,
			0x020F: 0x083C000000000000,
			0x0210: 0x0840000000000000,
			0x0211: 0x0844000000000000,
			0x0212: 0x184A124CECE3E120,
			0x0213: 0x184E132A27182130,
			0x0214: 0x185212A4D94C2140,
			0x0215: 0x6854000000000000,
			0x0216: 0x0858000000000000,
			0x0217: 0x185E0CBC92922170,
			0x0218: 0x0860000000000000,
			0x0219: 0x186612CCB0E26190,
			0x021A: 0x6868000000000000,
			0x0221: 0x188621A4F5182210,
			0x0222: 0x0888000000000000,
			0x0223: 0x088C000000000000,
			0x0224: 0x0890000000000000,
			0x0225: 0x0894000000000000,
			0x0226: 0x0898000000000000,
			0x0227: 0x089C000000000000,
			0x0228: 0x08A0000000000000,
			0x0229: 0x08A4000000000000,
			0x022A: 0x08A8000000000000,
			0x022B: 0x08AC000000000000,
			0x022C: 0x08B0000000000000,
			0x022D: 0x28B4000000000000,
			0x022E: 0x28B8000000000000,
			0x022F: 0x28BC000000000000,
			0x0230: 0x28C0000000000000,
			0x0231: 0x28C4000000000000,
			0x0232: 0x28C8000000000000,
			0x0233: 0x28CC000000000000,
			0x0234: 0x28D0000000000000,
			0x0235: 0x28D4000000000000,
			0x0236: 0x28D8000000000000,
			0x0237: 0x28DC000000000000,
			0x0238: 0x28E0000000000000,
			0x0239: 0x28E4000000000000,
			0x023A: 0x28E8000000000000,
			0x0241: 0x5904000000000000,
			0x0242: 0x5908000000000000,
			0x0243: 0x590C000000000000,
			0x0244: 0x0910000000000000,
			0x0245: 0x191645ACFC612450,
			0x0246: 0x0918000000000000,
			0x0247: 0x091C000000000000,
			0x0248: 0x5920000000000000,
			0x0249: 0x5924000000000000,
			0x024A: 0x0928000000000000,
			0x024B: 0x092C000000000000,
			0x024C: 0x5930000000000000,
			0x024D: 0x5934000000000000,
			0x024E: 0x5938000000000000,
			0x024F: 0x193E4FAD055A24F0,
			0x0250: 0x5940000000000000,
			0x0251: 0x0944000000000000,
			0x0252: 0x0948000000000000,
			0x0253: 0x194E7215614C2530,
			0x0254: 0x0950000000000000,
			0x0255: 0x1956559D075A2550,
			0x0256: 0x0958000000000000,
			0x0257: 0x195E570D0C612570,
			0x0258: 0x0960000000000000,
			0x0259: 0x0964000000000000,
			0x025A: 0x0968000000000000,
			0x0261: 0x198661AD55182610,
			0x0262: 0x198A6C10B4902620,
			0x0263: 0x198E791D64612630,
			0x0264: 0x19926475B2422640,
			0x0265: 0x1996772DE15C2650,
			0x0266: 0x5998000000000000,
			0x0267: 0x199E67857D0A2670,
			0x0268: 0x19A268751C60A680,
			0x0269: 0x19A67675834C2690,
			0x026A: 0x19AA6A6DD15C26A0,
			0x026B: 0x19AE765D7F5A26B0,
			0x026C: 0x19B26C2D6C60A6C0,
			0x026D: 0x19B66D95454C26D0,
			0x026E: 0x19BA65755C60A6E0,
			0x026F: 0x19BE6F6D846126F0,
			0x0270: 0x09C0000000000000,
			0x0271: 0x09C4000000000000,
			0x0272: 0x19CA7595C8E26720,
			0x0273: 0x19CE7325B0612730,
			0x0274: 0x19D274854C616740,
			0x0275: 0x39D4000000000000,
			0x0276: 0x19DA6CB1BCE2E760,
			0x0277: 0x09DC000000000000,
			0x0278: 0x19E2786C2CE3E780,
			0x0279: 0x19E67995F1182790,
			0x027A: 0x19EA77D5D861A7A0,
			0x0281: 0x3A04000000000000,
			0x0282: 0x0A08000000000000,
			0x0283: 0x1A0E830E38E3E830,
			0x0284: 0x1A12832128616840,
			0x0285: 0x0A14000000000000,
			0x0286: 0x1A18343208612860,
			0x0287: 0x1A1E877E0060A870,
			0x0288: 0x1A22880DF90A2880,
			0x0289: 0x0A24000000000000,
			0x028A: 0x1A2A8A5DF51628A0,
			0x028B: 0x1A2E8B66089428B0,
			0x028C: 0x1A328C9CE50A28C0,
			0x028D: 0x1A368B6E371628D0,
			0x028E: 0x1A3A9576284228E0,
			0x028F: 0x1A3E8F76109428F0,
			0x0290: 0x4A40000000000000,
			0x0291: 0x0A44000000000000,
			0x0292: 0x1A4A959631182920,
			0x0293: 0x0A4C000000000000,
			0x0294: 0x1A52947E18E3E940,
			0x0295: 0x0A54000000000000,
			0x0296: 0x1A5A95B63C942960,
			0x0297: 0x1A5E97713C002970,
			0x0298: 0x0A60000000000000,
			0x0299: 0x0A64000000000000,
			0x029A: 0x1A6A9A0E846129A0,
			0x02A1: 0x1A86AB96495A2A10,
			0x02A2: 0x0A88000000000000,
			0x02A3: 0x0A8C000000000000,
			0x02A4: 0x0A90000000000000,
			0x02A5: 0x0A94000000000000,
			0x02A6: 0x0A98000000000000,
			0x02A7: 0x1A9EA70E40612A70,
			0x02A8: 0x0AA0000000000000,
			0x02A9: 0x0AA4000000000000,
			0x02AA: 0x0AA8000000000000,
			0x02AB: 0x3AAC000000000000,
			0x02AC: 0x0AB0000000000000,
			0x02AD: 0x1AB6AD4C8A922B30,
			0x02AE: 0x3AB8000000000000,
			0x02AF: 0x0ABC000000000000,
			0x02B0: 0x0AC0000000000000,
			0x02B1: 0x0AC4000000000000,
			0x02B2: 0x0AC8000000000000,
			0x02B3: 0x1ACEB30E90C42B30,
			0x02B4: 0x0AD0000000000000,
			0x02B5: 0x0AD4000000000000,
			0x02B6: 0x0AD8000000000000,
			0x02B7: 0x0ADC000000000000,
			0x02B8: 0x0AE0000000000000,
			0x02B9: 0x1AE6B2CEB4E26B90,
			0x02BA: 0x1AEABA16B9162BA0,
			0x02C1: 0x1B06C1A2A14C2C10,
			0x02C2: 0x0B08000000000000,
			0x02C3: 0x1B0EC3A53CE3EC30,
			0x02C4: 0x6B10000000000000,
			0x02C5: 0x1B16C576BCE26C50,
			0x02C6: 0x0B18000000000000,
			0x02C7: 0x1B1EC710B8E3EC70,
			0x02C8: 0x0B20000000000000,
			0x02C9: 0x1B26C996A4E3EC90,
			0x02CA: 0x0B28000000000000,
			0x02CB: 0x0B2C000000000000,
			0x02CC: 0x0B30000000000000,
			0x02CD: 0x0B34000000000000,
			0x02CE: 0x1B3ACE6D810A2CE0,
			0x02CF: 0x0B3C000000000000,
			0x02D0: 0x0B40000000000000,
			0x02D1: 0x0B44000000000000,
			0x02D2: 0x0B48000000000000,
			0x02D3: 0x0B4C000000000000,
			0x02D4: 0x0B50000000000000,
			0x02D5: 0x1B56D5A448902D50,
			0x02D6: 0x0B58000000000000,
			0x02D7: 0x0B5C000000000000,
			0x02D8: 0x0B60000000000000,
			0x02D9: 0x0B64000000000000,
			0x02DA: 0x0B68000000000000,
			0x02E1: 0x0B84000000000000,
			0x02E2: 0x0B88000000000000,
			0x02E3: 0x0B8C000000000000,
			0x02E4: 0x0B90000000000000,
			0x02E5: 0x0B94000000000000,
			0x02E6: 0x1B9AEC36D8942E60,
			0x02E7: 0x5B9C000000000000,
			0x02E8: 0x0BA0000000000000,
			0x02E9: 0x0BA4000000000000,
			0x02EA: 0x0BA8000000000000,
			0x02EB: 0x6BAC000000000000,
			0x02EC: 0x5BB0000000000000,
			0x02ED: 0x0BB4000000000000,
			0x02EE: 0x0BB8000000000000,
			0x02EF: 0x5BBC000000000000,
			0x02F0: 0x0BC0000000000000,
			0x02F1: 0x0BC4000000000000,
			0x02F2: 0x0BC8000000000000,
			0x02F3: 0x1BCEF36EE4942F30,
			0x02F4: 0x0BD0000000000000,
			0x02F5: 0x0BD4000000000000,
			0x02F6: 0x5BD8000000000000,
			0x02F7: 0x0BDC000000000000,
			0x02F8: 0x0BE0000000000000,
			0x02F9: 0x0BE4000000000000,
			0x02FA: 0x0BE8000000000000,
			0x0301: 0x2C04000000000000,
			0x0302: 0x2C08000000000000,
			0x0303: 0x2C0C000000000000,
			0x0304: 0x2C10000000000000,
			0x0305: 0x2C14000000000000,
			0x0306: 0x2C18000000000000,
			0x0307: 0x2C1C000000000000,
			0x0308: 0x2C20000000000000,
			0x0309: 0x2C24000000000000,
			0x030A: 0x2C28000000000000,
			0x030B: 0x2C2C0000000016F0,
			0x030C: 0x2C30000000000000,
			0x030D: 0x2C34000000000000,
			0x030E: 0x2C38000000000000,
			0x030F: 0x2C3C000000000000,
			0x0310: 0x2C40000000000000,
			0x0311: 0x2C44000000000000,
			0x0312: 0x2C48000000000000,
			0x0313: 0x2C4C000000000000,
			0x0314: 0x2C50000000000000,
			0x0315: 0x2C54000000000000,
			0x0316: 0x2C58000000000000,
			0x0317: 0x2C5C000000000000,
			0x0318: 0x2C60000000000000,
			0x0319: 0x2C64000000000000,
			0x031A: 0x2C68000000000000,
			0x0321: 0x0C84000000000000,
			0x0322: 0x0C88000000000000,
			0x0323: 0x0C8C000000000000,
			0x0324: 0x6C90000000000000,
			0x0325: 0x1C97256EEF183250,
			0x0326: 0x0C98000000000000,
			0x0327: 0x0C9C000000000000,
			0x0328: 0x0CA0000000000000,
			0x0329: 0x0CA4000000000000,
			0x032A: 0x0CA8000000000000,
			0x032B: 0x0CAC000000000000,
			0x032C: 0x0CB0000000000000,
			0x032D: 0x0CB4000000000000,
			0x032E: 0x0CB8000000000000,
			0x032F: 0x0CBC000000000000,
			0x0330: 0x0CC0000000000000,
			0x0331: 0x0CC4000000000000,
			0x0332: 0x0CC8000000000000,
			0x0333: 0x0CCC000000000000,
			0x0334: 0x1CD1B9A15E613340,
			0x0335: 0x4CD4000000000000,
			0x0336: 0x5CD8000000000000,
			0x0337: 0x0CDC000000000000,
			0x0338: 0x0CE0000000000000,
			0x0339: 0x0CE4000000000000,
			0x033A: 0x0CE8000000000000,
			0x0341: 0x1D0741358C61B410,
			0x0342: 0x0D08000000000000,
			0x0343: 0x0D0C000000000000,
			0x0344: 0x0D10000000000000,
			0x0345: 0x0D14000000000000,
			0x0346: 0x0D18000000000000,
			0x0347: 0x0D1C000000000000,
			0x0348: 0x0D20000000000000,
			0x0349: 0x0D24000000000000,
			0x034A: 0x0D28000000000000,
			0x034B: 0x0D2C000000000000,
			0x034C: 0x0D30000000000000,
			0x034D: 0x1D374D16FC6134D0,
			0x034E: 0x0D38000000000000,
			0x034F: 0x0D3C000000000000,
			0x0350: 0x0D40000000000000,
			0x0351: 0x0D44000000000000,
			0x0352: 0x4D48000000000000,
			0x0353: 0x0D4C000000000000,
			0x0354: 0x0D50000000000000,
			0x0355: 0x0D54000000000000,
			0x0356: 0x0D58000000000000,
			0x0357: 0x1D5F572D98613570,
			0x0358: 0x0D60000000000000,
			0x0359: 0x0D64000000000000,
			0x035A: 0x2D68000000000000,
		},
	}

	// reverse indexes

	d.dStatusIDByName = make(map[string]uint8, len(d.dStatusByID))

	for k, v := range d.dStatusByID {
		d.dStatusIDByName[strings.ToUpper(v.name)] = uint8(k)
	}

	d.dRegionIDByCode = make(map[string]uint8, len(d.dRegionByID))
	d.dRegionIDByName = make(map[string]uint8, len(d.dRegionByID))

	for k, v := range d.dRegionByID {
		d.dRegionIDByCode[v.code] = uint8(k)
		d.dRegionIDByName[strings.ToUpper(v.name)] = uint8(k)
	}

	d.dSubRegionIDByCode = make(map[string]uint8, len(d.dSubRegionByID))
	d.dSubRegionIDByName = make(map[string]uint8, len(d.dSubRegionByID))

	for k, v := range d.dSubRegionByID {
		d.dSubRegionIDByCode[v.code] = uint8(k)
		d.dSubRegionIDByName[strings.ToUpper(v.name)] = uint8(k)
	}

	d.dIntermediateRegionIDByCode = make(map[string]uint8, len(d.dIntermediateRegionByID))
	d.dIntermediateRegionIDByName = make(map[string]uint8, len(d.dIntermediateRegionByID))

	for k, v := range d.dIntermediateRegionByID {
		d.dIntermediateRegionIDByCode[v.code] = uint8(k)
		d.dIntermediateRegionIDByName[strings.ToUpper(v.name)] = uint8(k)
	}

	// extra indexes

	d.dAlpha2IDByAlpha3ID = make(map[uint16]uint16, len(d.dCountryKeyByAlpha2ID))
	d.dAlpha2IDByNumericID = make(map[uint16]uint16, len(d.dCountryKeyByAlpha2ID))
	d.dAlpha2IDsByRegionID = make(map[uint8][]uint16, len(d.dRegionByID))
	d.dAlpha2IDsBySubRegionID = make(map[uint8][]uint16, len(d.dSubRegionByID))
	d.dAlpha2IDsByIntermediateRegionID = make(map[uint8][]uint16, len(d.dIntermediateRegionByID))
	d.dAlpha2IDsByStatusID = make(map[uint8][]uint16, len(d.dStatusByID))
	d.dAlpha2IDsByTLD = make(map[uint16][]uint16, len(d.dCountryKeyByAlpha2ID))

	for k, v := range d.dCountryKeyByAlpha2ID {
		ck := decodeCountryKey(v)

		d.dAlpha2IDByAlpha3ID[ck.alpha3] = k
		d.dAlpha2IDByNumericID[ck.numeric] = k
		d.dAlpha2IDsByRegionID[ck.region] = append(d.dAlpha2IDsByRegionID[ck.region], k)
		d.dAlpha2IDsBySubRegionID[ck.subregion] = append(d.dAlpha2IDsBySubRegionID[ck.subregion], k)
		d.dAlpha2IDsByIntermediateRegionID[ck.intregion] = append(d.dAlpha2IDsByIntermediateRegionID[ck.intregion], k)
		d.dAlpha2IDsByStatusID[ck.status] = append(d.dAlpha2IDsByStatusID[ck.status], k)
		d.dAlpha2IDsByTLD[ck.tld] = append(d.dAlpha2IDsByTLD[ck.tld], k)
	}

	delete(d.dAlpha2IDByAlpha3ID, 0)
	delete(d.dAlpha2IDByNumericID, 0)

	return d
}

func (d *Data) statusByID(id int) (*enumData, error) {
	if id < 0 || id >= len(d.dStatusByID) {
		return nil, errInvalidKey
	}

	return d.dStatusByID[id], nil
}

func (d *Data) statusIDByName(name string) (uint8, error) {
	v, ok := d.dStatusIDByName[strings.ToUpper(name)]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) regionByID(id int) (*enumData, error) {
	if id < 0 || id >= len(d.dRegionByID) {
		return nil, errInvalidKey
	}

	return d.dRegionByID[id], nil
}

func (d *Data) regionIDByCode(code string) (uint8, error) {
	v, ok := d.dRegionIDByCode[code]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) regionIDByName(name string) (uint8, error) {
	v, ok := d.dRegionIDByName[strings.ToUpper(name)]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) subRegionByID(id int) (*enumData, error) {
	if id < 0 || id >= len(d.dSubRegionByID) {
		return nil, errInvalidKey
	}

	return d.dSubRegionByID[id], nil
}

func (d *Data) subRegionIDByCode(code string) (uint8, error) {
	v, ok := d.dSubRegionIDByCode[code]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) subRegionIDByName(name string) (uint8, error) {
	v, ok := d.dSubRegionIDByName[strings.ToUpper(name)]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) intermediateRegionByID(id int) (*enumData, error) {
	if id < 0 || id >= len(d.dIntermediateRegionByID) {
		return nil, errInvalidKey
	}

	return d.dIntermediateRegionByID[id], nil
}

func (d *Data) intermediateRegionIDByCode(code string) (uint8, error) {
	v, ok := d.dIntermediateRegionIDByCode[code]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) intermediateRegionIDByName(name string) (uint8, error) {
	v, ok := d.dIntermediateRegionIDByName[strings.ToUpper(name)]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) countryNamesByAlpha2ID(id uint16) (*names, error) {
	v, ok := d.dCountryNamesByAlpha2ID[id]
	if !ok {
		return nil, errInvalidKey
	}

	return v, nil
}

func (d *Data) countryKeyByAlpha2ID(id uint16) (uint64, error) {
	v, ok := d.dCountryKeyByAlpha2ID[id]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) alpha2IDByAlpha3ID(id uint16) (uint16, error) {
	v, ok := d.dAlpha2IDByAlpha3ID[id]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) alpha2IDByNumericID(id uint16) (uint16, error) {
	v, ok := d.dAlpha2IDByNumericID[id]
	if !ok {
		return 0, errInvalidKey
	}

	return v, nil
}

func (d *Data) alpha2IDsByRegionID(id uint8) ([]uint16, error) {
	v, ok := d.dAlpha2IDsByRegionID[id]
	if !ok {
		return nil, errInvalidKey
	}

	return v, nil
}

func (d *Data) alpha2IDsBySubRegionID(id uint8) ([]uint16, error) {
	v, ok := d.dAlpha2IDsBySubRegionID[id]
	if !ok {
		return nil, errInvalidKey
	}

	return v, nil
}

func (d *Data) alpha2IDsByIntermediateRegionID(id uint8) ([]uint16, error) {
	v, ok := d.dAlpha2IDsByIntermediateRegionID[id]
	if !ok {
		return nil, errInvalidKey
	}

	return v, nil
}

func (d *Data) alpha2IDsByStatusID(id uint8) ([]uint16, error) {
	v, ok := d.dAlpha2IDsByStatusID[id]
	if !ok {
		return nil, errInvalidKey
	}

	return v, nil
}

func (d *Data) alpha2IDsByTLD(id uint16) ([]uint16, error) {
	v, ok := d.dAlpha2IDsByTLD[id]
	if !ok {
		return nil, errInvalidKey
	}

	return v, nil
}
