package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

var (
	ErrEmptyInput = errors.New("empty input")
)

func ConvertText(input string) (string, error) {
	if input == "" {
		return "", ErrEmptyInput
	}

	input = strings.TrimSpace(input)

	if isMorseCode(input) {
		result := morse.ToText(input)
		return result, nil
	}
	result := morse.ToMorse(input)
	return result, nil
}

func ConvertMorse(input string) (string, error) {
	if input == "" {
		return "", ErrEmptyInput
	}

	input = strings.TrimSpace(input)

	if !isMorseCode(input) {
		return "", errors.New("invalid morse code input")
	}

	result := morse.ToText(input)
	return result, nil
}

func isMorseCode(input string) bool {
	cleanInput := strings.ReplaceAll(input, " ", "")

	for _, char := range cleanInput {
		if char != '.' && char != '-' {
			return false
		}
	}

	return true
}
