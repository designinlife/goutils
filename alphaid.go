package goutils

// See https://github.com/dineshappavoo/basex

import (
	"errors"
	"math/big"
	"strconv"
	"unicode"
)

var (
	dictionary = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	base       *big.Int
	dictMap    map[byte]*big.Int
)

func init() {
	base = big.NewInt(int64(len(dictionary)))

	// for efficiency, make a map
	dictMap = make(map[byte]*big.Int)

	j := 0
	for _, val := range dictionary {
		dictMap[val] = big.NewInt(int64(j))
		j = j + 1
	}
}

// checks if given string is a valid numeric
func isValidNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

// checks if s is ascii and printable, aka doesn't include tab, backspace, etc.
func isAsciiPrintable(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) || unicode.IsPunct(r) {
			return false
		}
	}
	return true
}

// encodeInt encodes a big.Int integer, the value of remaining is changed to 0 during the process
func encodeInt(remaining *big.Int) (string, error) {
	var result []byte
	var index int
	var strVal string

	a := big.NewInt(0)
	b := big.NewInt(0)
	c := big.NewInt(0)
	d := big.NewInt(0)

	exponent := 1

	for remaining.Cmp(big.NewInt(0)) != 0 {
		a.Exp(base, big.NewInt(int64(exponent)), nil) // 16^1 = 16
		b = b.Mod(remaining, a)                       // 119 % 16 = 7 | 112 % 256 = 112
		c = c.Exp(base, big.NewInt(int64(exponent-1)), nil)
		d = d.Div(b, c)

		// if d > dictionary.length, we have a problem. but BigInteger doesn't have
		// a greater than method :-(  hope for the best. theoretically, d is always
		// an index of the dictionary!
		strVal = d.String()
		index, _ = strconv.Atoi(strVal)
		result = append(result, dictionary[index])
		remaining = remaining.Sub(remaining, b) // 119 - 7 = 112 | 112 - 112 = 0
		exponent = exponent + 1
	}

	// need to reverse it, since the start of the list contains the least significant values
	return string(reverse(result)), nil
}

// EncodeInt encodes a big.Int integer
func EncodeInt(i *big.Int) (string, error) {
	remaining := big.NewInt(0)
	remaining.Set(i)

	return encodeInt(remaining)
}

// Encode converts the big integer to alpha id (an alphanumeric id with mixed cases)
func Encode(s string) (string, error) {
	// numeric validation
	if !isValidNumeric(s) {
		return "", errors.New("encode string is not a valid numeric")
	}

	remaining := big.NewInt(0)
	remaining.SetString(s, 10)

	return encodeInt(remaining)
}

// DecodeInt converts the alpha id to a bit.Int
func DecodeInt(s string) (*big.Int, error) {
	// Validate if given string is valid
	if !isAsciiPrintable(s) {
		return nil, errors.New("decode string is not valid.[a-z, A_Z, 0-9] only allowed")
	}
	// reverse it, coz its already reversed!
	chars2 := reverse([]byte(s))

	bi := big.NewInt(0)

	exponent := 0
	a := big.NewInt(0)
	b := big.NewInt(0)
	intermed := big.NewInt(0)

	for _, c := range chars2 {
		a = dictMap[c]
		intermed = intermed.Exp(base, big.NewInt(int64(exponent)), nil)
		b = b.Mul(intermed, a)
		bi = bi.Add(bi, b)
		exponent = exponent + 1
	}
	return bi, nil
}

// Decode converts the alpha id to big integer
func Decode(s string) (string, error) {
	bi, err := DecodeInt(s)
	if err != nil {
		return "", err
	}
	return bi.String(), nil
}

func reverse(bs []byte) []byte {
	for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
		bs[i], bs[j] = bs[j], bs[i]
	}
	return bs
}
