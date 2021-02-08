package server

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const idLength = 8

func newID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	alphabetRunes := []rune(alphabet)

	runes := make([]rune, idLength)
	for i := 0; i < idLength; i++ {
		runes[i] = alphabetRunes[rand.Intn(len(alphabetRunes))]
	}

	return string(runes)
}
