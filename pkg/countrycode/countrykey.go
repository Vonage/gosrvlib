package countrycode

import (
	"errors"
	"fmt"
	"strconv"
)

// Binary bit masks for when each countrykey element starts and ends.
const (
	bitMaskStatus    uint64 = 0x7000000000000000 // 01110000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
	bitMaskAlpha2    uint64 = 0x0FFC000000000000 // 00001111 11111100 00000000 00000000 00000000 00000000 00000000 00000000
	bitMaskAlpha3    uint64 = 0x0003FFF800000000 // 00000000 00000011 11111111 11111000 00000000 00000000 00000000 00000000
	bitMaskNumeric   uint64 = 0x00000007FE000000 // 00000000 00000000 00000000 00000111 11111110 00000000 00000000 00000000
	bitMaskRegion    uint64 = 0x0000000001C00000 // 00000000 00000000 00000000 00000000 00000001 11000000 00000000 00000000
	bitMaskSubRegion uint64 = 0x00000000003E0000 // 00000000 00000000 00000000 00000000 00000000 00111110 00000000 00000000
	bitMaskIntRegion uint64 = 0x000000000001C000 // 00000000 00000000 00000000 00000000 00000000 00000001 11000000 00000000
	bitMaskTLD       uint64 = 0x0000000000003FF0 // 00000000 00000000 00000000 00000000 00000000 00000000 00111111 11110000
)

// Binary bit positions for when each countrykey element starts (counting from the right - LSB).
const (
	bitPosStatus    int = 60
	bitPosAlpha2    int = 50
	bitPosAlpha3    int = 35
	bitPosNumeric   int = 25
	bitPosRegion    int = 22
	bitPosSubRegion int = 17
	bitPosIntRegion int = 14
	bitPosTLD       int = 4
)

// Binary bit masks for when each character in the alpha-2, alpha-3 and TLD codes starts and ends.
const (
	bitMaskChar2 uint16 = 0x7C00 // 111110000000000
	bitMaskChar1 uint16 = 0x03E0 // 000001111100000
	bitMaskChar0 uint16 = 0x001F // 000000000011111
)

// Binary bit positions for when each character in the alpha-2, alpha-3 and TLD codes starts (counting from the right - LSB).
const (
	bitPosChar1 int = 5
	bitPosChar2 int = 10
)

const (
	chrOffsetUpper uint16 = 64 // the character 'A' is 65 in ASCII, so we need to shift it to 1
	chrOffsetLower uint16 = 96 // the character 'a' is 97 in ASCII, so we need to shift it to 1
)

var (
	errInvalidKey       = errors.New("invalid key")
	errInvalidLength    = errors.New("invalid code length")
	errInvalidCharacter = errors.New("invalid code, it must contain only A-Z characters")
)

// countryKeyElem represent a CountryKey uint64 format:
//
//	-----------------------
//	 1 bit Reserved
//	 3 bit Status
//	10 bit Alpha-2 Code
//	15 bit Alpha-3 Code
//	10 bit Numeric Code
//	 3 bit Region Enum
//	 5 bit Sub-Region Enum
//	 3 bit Intermediate-Region Enum
//	10 bit TLD
//	 4 bit Reserved
//	-----------------------
type countryKeyElem struct {
	status    uint8  // 3 bit status
	alpha2    uint16 // 10 bit alpha2
	alpha3    uint16 // 15 bit alpha3
	numeric   uint16 // 10 bit numeric
	region    uint8  // 3 bit region
	subregion uint8  // 5 bit sub-region
	intregion uint8  // 3 bit intermediate-region
	tld       uint16 // 10 bit TLD
}

func decodeCountryKey(key uint64) *countryKeyElem {
	return &countryKeyElem{
		status:    uint8((key & bitMaskStatus) >> bitPosStatus),
		alpha2:    uint16((key & bitMaskAlpha2) >> bitPosAlpha2),
		alpha3:    uint16((key & bitMaskAlpha3) >> bitPosAlpha3),
		numeric:   uint16((key & bitMaskNumeric) >> bitPosNumeric),
		region:    uint8((key & bitMaskRegion) >> bitPosRegion),
		subregion: uint8((key & bitMaskSubRegion) >> bitPosSubRegion),
		intregion: uint8((key & bitMaskIntRegion) >> bitPosIntRegion),
		tld:       uint16((key & bitMaskTLD) >> bitPosTLD),
	}
}

func charOffset(b byte, offset uint16) (uint16, error) {
	c := (uint16(b) - offset)
	if c < 1 || c > 26 { // A-Z or a-z
		return 0, errInvalidCharacter
	}

	return c, nil
}

func charOffsetUpper(b byte) (uint16, error) {
	return charOffset(b, chrOffsetUpper)
}

func charOffsetLower(b byte) (uint16, error) {
	return charOffset(b, chrOffsetLower)
}

func encodeAlpha2(s string) (uint16, error) {
	if len(s) != 2 {
		return 0, errInvalidLength
	}

	c0, err := charOffsetUpper(s[0])
	if err != nil {
		return c0, err
	}

	c1, err := charOffsetUpper(s[1])
	if err != nil {
		return c1, err
	}

	return (c0<<5 | c1), nil
}

func decodeAlpha2(code uint16) string {
	return string([]byte{
		byte(((code & bitMaskChar1) >> bitPosChar1) + chrOffsetUpper),
		byte((code & bitMaskChar0) + chrOffsetUpper),
	})
}

func encodeAlpha3(s string) (uint16, error) {
	if len(s) != 3 {
		return 0, errInvalidLength
	}

	c0, err := charOffsetUpper(s[0])
	if err != nil {
		return c0, err
	}

	c1, err := charOffsetUpper(s[1])
	if err != nil {
		return c1, err
	}

	c2, err := charOffsetUpper(s[2])
	if err != nil {
		return c2, err
	}

	return (c0<<10 | c1<<5 | c2), nil
}

func decodeAlpha3(code uint16) string {
	return string([]byte{
		byte(((code & bitMaskChar2) >> bitPosChar2) + chrOffsetUpper),
		byte(((code & bitMaskChar1) >> bitPosChar1) + chrOffsetUpper),
		byte((code & bitMaskChar0) + chrOffsetUpper),
	})
}

func encodeTLD(s string) (uint16, error) {
	if len(s) != 2 {
		return 0, errInvalidLength
	}

	c0, err := charOffsetLower(s[0])
	if err != nil {
		return c0, err
	}

	c1, err := charOffsetLower(s[1])
	if err != nil {
		return c1, err
	}

	return (c0<<5 | c1), nil
}

func decodeTLD(code uint16) string {
	return string([]byte{
		byte(((code & bitMaskChar1) >> bitPosChar1) + chrOffsetLower),
		byte((code & bitMaskChar0) + chrOffsetLower),
	})
}

func encodeNumeric(s string) (uint16, error) {
	if len(s) != 3 {
		return 0, errInvalidLength
	}

	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric code: %w", err)
	}

	return uint16(v), nil
}
