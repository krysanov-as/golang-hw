package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	runes := []rune(str)
	if len(runes) == 0 {
		return "", nil
	}

	var result strings.Builder
	var prevRune rune

	for _, r := range runes {
		if unicode.IsDigit(r) {
			if prevRune == 0 {
				return "", ErrInvalidString
			}

			n, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			if n == 0 {
				prevRune = 0
				continue
			}

			result.WriteString(strings.Repeat(string(prevRune), n))
			prevRune = 0
			continue
		}

		if prevRune != 0 {
			result.WriteRune(prevRune)
		}
		prevRune = r
	}

	if prevRune != 0 {
		result.WriteRune(prevRune)
	}

	return result.String(), nil
}
