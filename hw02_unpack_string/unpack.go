package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result strings.Builder
	var cursorAlhpa rune
	var escapedNext bool
	unpack := func(sym rune, i int) {
		result.WriteString(strings.Repeat(string(cursorAlhpa), i))
		cursorAlhpa = 0
		escapedNext = false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return "", ErrInvalidString
		}
		if !unicode.IsDigit(r) {
			// check for backslash
			if r == 0x5C {
				if cursorAlhpa != 0 {
					unpack(cursorAlhpa, 1)
				}
				if !escapedNext {
					escapedNext = true
				} else {
					cursorAlhpa = r
					escapedNext = false
				}

				continue
			} else {
				if cursorAlhpa != 0 {
					unpack(cursorAlhpa, 1)
				}
				cursorAlhpa = r
			}
		} else {
			//is digit
			if escapedNext {
				cursorAlhpa = r
				escapedNext = false
				continue
				//
			} else {
				if cursorAlhpa == 0 {
					return "", ErrInvalidString
				}
				digit, err := strconv.Atoi(string(r))
				if err != nil {
					return "", err
				}
				unpack(cursorAlhpa, digit)
			}
			if escapedNext {
				cursorAlhpa = r
			}
		}
	}
	// check last symbol
	if cursorAlhpa != 0 {
		unpack(cursorAlhpa, 1)
	}

	return result.String(), nil
}
