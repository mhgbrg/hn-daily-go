package models

import "math/rand"

const userIDLength = 6
const userIDChars = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateUserID() string {
	b := make([]byte, userIDLength)
	for i := range b {
		b[i] = userIDChars[rand.Intn(len(userIDChars))]
	}
	return string(b)
}
