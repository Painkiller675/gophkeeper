package utils

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano()) // TODO: use unpredicted generator here
}

// errInvalidLength is an invalid length error
var errInvalidLength = errors.New("invalid length")

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString generates random string of the set length генерирует случайную строку заданной длины из символов letters
func RandomString(length int) (string, error) {
	if length <= 0 {
		return "", errInvalidLength
	}
	// use builder to concatenate
	var builder strings.Builder
	builder.Grow(length)
	for i := 0; i < length; i++ {
		index := rand.Intn(len(letters))
		builder.WriteRune(letters[index])
	}
	return builder.String(), nil
}
