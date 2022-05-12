package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	sRunes := []rune(s)
	if unicode.IsDigit(sRunes[0]) {
		return "", ErrInvalidString
	}

	var b strings.Builder
	for i, rCount := 0, len(sRunes); i < rCount; i++ {
		r := sRunes[i]

		if r == '\\' {
			if i == rCount-1 {
				return "", ErrInvalidString
			}

			if i < rCount-2 && unicode.IsDigit(sRunes[i+2]) {
				tmp, _ := strconv.Atoi(string(sRunes[i+2]))
				b.WriteString(strings.Repeat(string(sRunes[i+1]), tmp))
				i += 2

				continue
			}
			b.WriteRune(sRunes[i+1])
			i++

			continue
		}

		if unicode.IsDigit(r) {
			if i < rCount-1 && unicode.IsDigit(sRunes[i+1]) {
				return "", ErrInvalidString
			}

			continue
		}

		if i < rCount-1 && unicode.IsDigit(sRunes[i+1]) {
			tmp, _ := strconv.Atoi(string(sRunes[i+1]))
			b.WriteString(strings.Repeat(string(r), tmp))

			continue
		}

		b.WriteRune(r)
	}

	return b.String(), nil
}
