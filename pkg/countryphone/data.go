package countryphone

// defaultData returns a default map of phone numbers to country ISO-3166 Alpha-2 Codes.
// Ref.: https://en.wikipedia.org/wiki/List_of_country_calling_codes
//
//nolint:funlen,maintidx
func defaultData() InData {
	return InData{
		"AC": {CC: "247"},
		"AD": {CC: "376"},
		"AE": {CC: "971"},
		"AF": {CC: "93"},
		"AG": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1268"}},
		}},
		"AI": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1264"}},
		}},
		"AL": {CC: "355"},
		"AM": {CC: "374"},
		"AO": {CC: "244"},
		"AQ": {CC: "672"},
		"AR": {CC: "54"},
		"AS": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1684"}},
		}},
		"AT": {CC: "43"},
		"AU": {CC: "61"},
		"AW": {CC: "297"},
		"AX": {CC: "35818"},
		"AZ": {CC: "994"},
		"BA": {CC: "387"},
		"BB": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1246"}},
		}},
		"BD": {CC: "880"},
		"BE": {CC: "32"},
		"BF": {CC: "226"},
		"BG": {CC: "359"},
		"BH": {CC: "973"},
		"BI": {CC: "257"},
		"BJ": {CC: "229"},
		"BL": {CC: "590"},
		"BM": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1441"}},
		}},
		"BN": {CC: "673"},
		"BO": {CC: "591"},
		"BQ": {CC: "599"},
		"BR": {CC: "55"},
		"BS": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1242"}},
		}},
		"BT": {CC: "975"},
		"BV": {CC: "47"},
		"BW": {CC: "267"},
		"BY": {CC: "375"},
		"BZ": {CC: "501"},
		"CA": {CC: "1", Groups: []InPrefixGroup{
			{"Alberta", 2, 1, []string{"1368", "1403", "1568", "1587", "1780", "1825"}},
			{"British Columbia", 2, 1, []string{"1236", "1250", "1257", "1604", "1672", "1778"}},
			{"Manitoba", 2, 1, []string{"1204", "1431", "1584"}},
			{"New Brunswick", 2, 1, []string{"1428", "1506"}},
			{"Newfoundland and Labrador", 2, 1, []string{"1709", "1879"}},
			{"Northwest Territories", 2, 1, []string{"1867"}},
			{"Nova Scotia", 2, 1, []string{"1782", "1851", "1902"}},
			{"Nunavut", 2, 1, []string{"1867"}},
			{"Ontario", 2, 1, []string{"1226", "1249", "1289", "1343", "1365", "1382", "1387", "1416", "1437", "1460", "1519", "1537", "1548", "1613", "1647", "1683", "1705", "1742", "1753", "1807", "1905", "1942"}},
			{"Prince Edward Island", 2, 1, []string{"1782", "1902"}},
			{"Quebec", 2, 1, []string{"1263", "1354", "1367", "1418", "1438", "1450", "1468", "1514", "1579", "1581", "1819", "1873"}},
			{"Saskatchewan", 2, 1, []string{"1306", "1474", "1639"}},
			{"Yukon", 2, 1, []string{"1867"}},
			{"Canada special services", 4, 5, []string{"1600", "1622", "1633", "1644", "1655", "1677", "1688"}},
		}},
		"CC": {CC: "61"},
		"CD": {CC: "243"},
		"CF": {CC: "236"},
		"CG": {CC: "242"},
		"CH": {CC: "41"},
		"CI": {CC: "225"},
		"CK": {CC: "682"},
		"CL": {CC: "56"},
		"CM": {CC: "237"},
		"CN": {CC: "86"},
		"CO": {CC: "57"},
		"CR": {CC: "506"},
		"CT": {CC: "90"},
		"CU": {CC: "53"},
		"CV": {CC: "238"},
		"CW": {CC: "599"},
		"CX": {CC: "61"},
		"CY": {CC: "357"},
		"CZ": {CC: "420"},
		"DE": {CC: "49"},
		"DJ": {CC: "253"},
		"DK": {CC: "45"},
		"DM": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1767"}},
		}},
		"DO": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1809", "1829", "1849"}},
		}},
		"DZ": {CC: "213"},
		"EC": {CC: "593"},
		"EE": {CC: "372"},
		"EG": {CC: "20"},
		"EH": {CC: "212"},
		"ER": {CC: "291"},
		"ES": {CC: "34"},
		"ET": {CC: "251"},
		"FI": {CC: "358"},
		"FJ": {CC: "679"},
		"FK": {CC: "500"},
		"FM": {CC: "691"},
		"FO": {CC: "298"},
		"FR": {CC: "33"},
		"GA": {CC: "241"},
		"GB": {CC: "44"},
		"GD": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1473"}},
		}},
		"GE": {CC: "995"},
		"GF": {CC: "594"},
		"GG": {CC: "44"},
		"GH": {CC: "233"},
		"GI": {CC: "350"},
		"GL": {CC: "299"},
		"GM": {CC: "220"},
		"GN": {CC: "224"},
		"GP": {CC: "590"},
		"GQ": {CC: "240"},
		"GR": {CC: "30"},
		"GS": {CC: "500"},
		"GT": {CC: "502"},
		"GU": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1671"}},
		}},
		"GW": {CC: "245"},
		"GY": {CC: "592"},
		"HK": {CC: "852"},
		"HN": {CC: "504"},
		"HR": {CC: "385"},
		"HT": {CC: "509"},
		"HU": {CC: "36"},
		"ID": {CC: "62"},
		"IE": {CC: "353"},
		"IL": {CC: "972"},
		"IM": {CC: "44"},
		"IN": {CC: "91"},
		"IO": {CC: "246"},
		"IQ": {CC: "964"},
		"IR": {CC: "98"},
		"IS": {CC: "354"},
		"IT": {CC: "39"},
		"JE": {CC: "44"},
		"JM": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1658", "1876"}},
		}},
		"JO": {CC: "962"},
		"JP": {CC: "81"},
		"KE": {CC: "254"},
		"KG": {CC: "996"},
		"KH": {CC: "855"},
		"KI": {CC: "686"},
		"KM": {CC: "269"},
		"KN": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1869"}},
		}},
		"KP": {CC: "850"},
		"KR": {CC: "82"},
		"KW": {CC: "965"},
		"KY": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1345"}},
		}},
		"KZ": {CC: "7", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"76", "77"}},
		}},
		"LA": {CC: "856"},
		"LB": {CC: "961"},
		"LC": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1758"}},
		}},
		"LI": {CC: "423"},
		"LK": {CC: "94"},
		"LR": {CC: "231"},
		"LS": {CC: "266"},
		"LT": {CC: "370"},
		"LU": {CC: "352"},
		"LV": {CC: "371"},
		"LY": {CC: "218"},
		"MA": {CC: "212"},
		"MC": {CC: "377"},
		"MD": {CC: "373"},
		"ME": {CC: "382"},
		"MF": {CC: "590"},
		"MG": {CC: "261"},
		"MH": {CC: "692"},
		"MK": {CC: "389"},
		"ML": {CC: "223"},
		"MM": {CC: "95"},
		"MN": {CC: "976"},
		"MO": {CC: "853"},
		"MP": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1670"}},
		}},
		"MQ": {CC: "596"},
		"MR": {CC: "222"},
		"MS": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1664"}},
		}},
		"MT": {CC: "356"},
		"MU": {CC: "230"},
		"MV": {CC: "960"},
		"MW": {CC: "265"},
		"MX": {CC: "52"},
		"MY": {CC: "60"},
		"MZ": {CC: "258"},
		"NA": {CC: "264"},
		"NC": {CC: "687"},
		"NE": {CC: "227"},
		"NF": {CC: "672"},
		"NG": {CC: "234"},
		"NI": {CC: "505"},
		"NL": {CC: "31"},
		"NO": {CC: "47"},
		"NP": {CC: "977"},
		"NR": {CC: "674"},
		"NU": {CC: "683"},
		"NZ": {CC: "64"},
		"OM": {CC: "968"},
		"PA": {CC: "507"},
		"PE": {CC: "51"},
		"PF": {CC: "689"},
		"PG": {CC: "675"},
		"PH": {CC: "63"},
		"PK": {CC: "92"},
		"PL": {CC: "48"},
		"PM": {CC: "508"},
		"PN": {CC: "64"},
		"PR": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1787", "1939"}},
		}},
		"PS": {CC: "970"},
		"PT": {CC: "351"},
		"PW": {CC: "680"},
		"PY": {CC: "595"},
		"QA": {CC: "974"},
		"QN": {CC: "374"},
		"RE": {CC: "262"},
		"RO": {CC: "40"},
		"RS": {CC: "381"},
		"RU": {CC: "7", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"71", "73", "74", "75", "78", "79"}},
		}},
		"RW": {CC: "250"},
		"SA": {CC: "966"},
		"SB": {CC: "677"},
		"SC": {CC: "248"},
		"SD": {CC: "249"},
		"SE": {CC: "46"},
		"SG": {CC: "65"},
		"SH": {CC: "290"},
		"SI": {CC: "386"},
		"SJ": {CC: "47"},
		"SK": {CC: "421"},
		"SL": {CC: "232"},
		"SM": {CC: "378"},
		"SN": {CC: "221"},
		"SO": {CC: "252"},
		"SR": {CC: "597"},
		"SS": {CC: "211"},
		"ST": {CC: "239"},
		"SV": {CC: "503"},
		"SX": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1721"}},
		}},
		"SY": {CC: "963"},
		"SZ": {CC: "268"},
		"TA": {CC: "290"},
		"TC": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1649"}},
		}},
		"TD": {CC: "235"},
		"TF": {CC: "262"},
		"TG": {CC: "228"},
		"TH": {CC: "66"},
		"TJ": {CC: "992"},
		"TK": {CC: "690"},
		"TL": {CC: "670"},
		"TM": {CC: "993"},
		"TN": {CC: "216"},
		"TO": {CC: "676"},
		"TR": {CC: "90"},
		"TT": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1868"}},
		}},
		"TV": {CC: "688"},
		"TW": {CC: "886"},
		"TZ": {CC: "255"},
		"UA": {CC: "380"},
		"UG": {CC: "256"},
		"UN": {CC: "888"},
		"US": {CC: "1", Groups: []InPrefixGroup{
			{"Alabama", 1, 1, []string{"1205", "1251", "1256", "1334", "1483", "1659", "1938"}},
			{"Alaska", 1, 1, []string{"1907"}},
			{"Arizona", 1, 1, []string{"1480", "1520", "1602", "1623", "1928"}},
			{"Arkansas", 1, 1, []string{"1327", "1479", "1501", "1870"}},
			{"California", 1, 1, []string{"1209", "1213", "1279", "1310", "1323", "1341", "1350", "1357", "1369", "1408", "1415", "1424", "1442", "1510", "1530", "1559", "1562", "1619", "1626", "1628", "1650", "1657", "1661", "1669", "1707", "1714", "1738", "1747", "1760", "1805", "1818", "1820", "1831", "1837", "1840", "1858", "1909", "1916", "1925", "1949", "1951"}},
			{"Colorado", 1, 1, []string{"1303", "1719", "1720", "1748", "1970", "1983"}},
			{"Connecticut", 1, 1, []string{"1203", "1475", "1860", "1959"}},
			{"Delaware", 1, 1, []string{"1302"}},
			{"District of Columbia", 1, 1, []string{"1202", "1771"}},
			{"Florida", 1, 1, []string{"1239", "1305", "1321", "1324", "1352", "1386", "1407", "1448", "1561", "1645", "1656", "1689", "1727", "1728", "1754", "1772", "1786", "1813", "1850", "1863", "1904", "1941", "1954"}},
			{"Georgia", 1, 1, []string{"1229", "1404", "1470", "1478", "1678", "1706", "1762", "1770", "1912", "1943"}},
			{"Hawaii", 1, 1, []string{"1808"}},
			{"Idaho", 1, 1, []string{"1208", "1986"}},
			{"Illinois", 1, 1, []string{"1217", "1224", "1309", "1312", "1331", "1447", "1464", "1618", "1630", "1708", "1730", "1773", "1779", "1815", "1847", "1861", "1872"}},
			{"Indiana", 1, 1, []string{"1219", "1260", "1317", "1463", "1574", "1765", "1812", "1930"}},
			{"Iowa", 1, 1, []string{"1319", "1515", "1563", "1641", "1712"}},
			{"Kansas", 1, 1, []string{"1316", "1620", "1785", "1913"}},
			{"Kentucky", 1, 1, []string{"1270", "1364", "1502", "1606", "1859"}},
			{"Louisiana", 1, 1, []string{"1225", "1318", "1337", "1457", "1504", "1985"}},
			{"Maine", 1, 1, []string{"1207"}},
			{"Maryland", 1, 1, []string{"1227", "1240", "1301", "1410", "1443", "1667"}},
			{"Massachusetts", 1, 1, []string{"1339", "1351", "1413", "1508", "1617", "1774", "1781", "1857", "1978"}},
			{"Michigan", 1, 1, []string{"1231", "1248", "1269", "1313", "1517", "1586", "1616", "1679", "1734", "1810", "1906", "1947", "1989"}},
			{"Minnesota", 1, 1, []string{"1218", "1320", "1507", "1612", "1651", "1763", "1924", "1952"}},
			{"Mississippi", 1, 1, []string{"1228", "1471", "1601", "1662", "1769"}},
			{"Missouri", 1, 1, []string{"1235", "1314", "1417", "1557", "1573", "1636", "1660", "1816", "1975"}},
			{"Montana", 1, 1, []string{"1406"}},
			{"Nebraska", 1, 1, []string{"1308", "1402", "1531"}},
			{"Nevada", 1, 1, []string{"1702", "1725", "1775"}},
			{"New Hampshire", 1, 1, []string{"1603"}},
			{"New Jersey", 1, 1, []string{"1201", "1551", "1609", "1640", "1732", "1848", "1856", "1862", "1908", "1973"}},
			{"New Mexico", 1, 1, []string{"1505", "1575"}},
			{"New York", 1, 1, []string{"1212", "1315", "1329", "1332", "1347", "1363", "1516", "1518", "1585", "1607", "1624", "1631", "1646", "1680", "1716", "1718", "1838", "1845", "1914", "1917", "1929", "1934"}},
			{"North Carolina", 1, 1, []string{"1252", "1336", "1472", "1704", "1743", "1828", "1910", "1919", "1980", "1984"}},
			{"North Dakota", 1, 1, []string{"1701"}},
			{"Ohio", 1, 1, []string{"1216", "1220", "1234", "1283", "1326", "1330", "1380", "1419", "1436", "1440", "1513", "1567", "1614", "1740", "1937"}},
			{"Oklahoma", 1, 1, []string{"1405", "1539", "1572", "1580", "1918"}},
			{"Oregon", 1, 1, []string{"1458", "1503", "1541", "1971"}},
			{"Pennsylvania", 1, 1, []string{"1215", "1223", "1267", "1272", "1412", "1445", "1484", "1570", "1582", "1610", "1717", "1724", "1814", "1835", "1878"}},
			{"Rhode Island", 1, 1, []string{"1401"}},
			{"South Carolina", 1, 1, []string{"1803", "1821", "1839", "1843", "1854", "1864"}},
			{"South Dakota", 1, 1, []string{"1605"}},
			{"Tennessee", 1, 1, []string{"1423", "1615", "1629", "1731", "1865", "1901", "1931"}},
			{"Texas", 1, 1, []string{"1210", "1214", "1254", "1281", "1325", "1346", "1361", "1409", "1430", "1432", "1469", "1512", "1621", "1682", "1713", "1726", "1737", "1806", "1817", "1830", "1832", "1903", "1915", "1936", "1940", "1945", "1956", "1972", "1979"}},
			{"Utah", 1, 1, []string{"1385", "1435", "1801"}},
			{"Vermont", 1, 1, []string{"1802"}},
			{"Virginia", 1, 1, []string{"1276", "1434", "1540", "1571", "1686", "1703", "1757", "1804", "1826", "1948"}},
			{"Washington", 1, 1, []string{"1206", "1253", "1360", "1425", "1509", "1564"}},
			{"West Virginia", 1, 1, []string{"1304", "1681"}},
			{"Wisconsin", 1, 1, []string{"1262", "1274", "1353", "1414", "1534", "1608", "1715", "1920"}},
			{"Wyoming", 1, 1, []string{"1307"}},
			{"Midway Atoll, Wake Island, Hawaii", 1, 1, []string{"1808"}},
			{"US government", 4, 5, []string{"1710"}},
		}},
		"UY": {CC: "598"},
		"UZ": {CC: "998"},
		"VA": {CC: "39", Groups: []InPrefixGroup{
			{"Vatican City", 0, 1, []string{"3906698"}},
			{"", 0, 0, []string{"379"}},
		}},
		"VC": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1784"}},
		}},
		"VE": {CC: "58"},
		"VG": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1284"}},
		}},
		"VI": {CC: "1", Groups: []InPrefixGroup{
			{"", 0, 0, []string{"1340"}},
		}},
		"VN": {CC: "84"},
		"VU": {CC: "678"},
		"WF": {CC: "681"},
		"WS": {CC: "685"},
		"XK": {CC: "383"},
		"YE": {CC: "967"},
		"YT": {CC: "262"},
		"ZA": {CC: "27"},
		"ZM": {CC: "260"},
		"ZW": {CC: "263"},
		"__": {CC: "", Groups: []InPrefixGroup{
			{"Universal International Freephone Service", 4, 5, []string{"800"}},
			{"Universal International Shared Cost Number", 4, 5, []string{"808"}},
			{"Inmarsat", 4, 5, []string{"870"}},
			{"Universal Personal Telecommunications", 4, 5, []string{"878"}},
			{"Global Mobile Satellite System", 4, 5, []string{"881"}},
			{"International Networks", 4, 5, []string{"882", "883"}},
			{"International premium rate service", 4, 5, []string{"979"}},
			{"International Telecommunications Public Correspondence Service", 4, 5, []string{"991"}},
		}},
	}
}
