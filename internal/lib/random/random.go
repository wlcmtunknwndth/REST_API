package random

import (
	"math/rand"
	"time"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewRandomString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var rnd []byte = make([]byte, length)
	for i := range rnd {
		rnd[i] = letters[rand.Intn(len(letters))]
	}
	return string(rnd)
}
