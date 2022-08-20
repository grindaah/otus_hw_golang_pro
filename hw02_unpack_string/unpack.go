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
	unpack := func(i int) {
		result.WriteString(strings.Repeat(string(cursorAlhpa), i))
		cursorAlhpa = 0
		escapedNext = false
	}
	escapeNext := func(sym rune) {
		if !escapedNext {
			escapedNext = true
		} else {
			cursorAlhpa = sym
			escapedNext = false
		}
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return "", ErrInvalidString
		}
		if r == 0x5C {
			if cursorAlhpa != 0 {
				unpack(1)
			}
			escapeNext(r)
			continue
		}
		if unicode.IsDigit(r) && !escapedNext && cursorAlhpa == 0 {
			return "", ErrInvalidString
		}
		if !unicode.IsDigit(r) {
			if cursorAlhpa != 0 {
				unpack(1)
			}
			cursorAlhpa = r
		} else {
			if escapedNext {
				cursorAlhpa = r
				escapedNext = false
			} else {
				digit, err := strconv.Atoi(string(r))
				if err != nil {
					return "", err
				}
				unpack(digit)
			}
		}
	}

	if cursorAlhpa != 0 {
		unpack(1)
	}

	return result.String(), nil
}
