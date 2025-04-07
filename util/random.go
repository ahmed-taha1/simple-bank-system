package util

import (
	"math/rand"
	"strings"
	"time"
)

var randGenerator *rand.Rand

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Create a new random generator
func init() {
	var randSource = rand.NewSource(time.Now().UnixNano())
	randGenerator = rand.New(randSource)
}

func RandomInt(min, max int64) int64 {
	return randGenerator.Int63n(max-min+1) + min // [min, max]
}

func RandomString(n int) string {
	var strBuilder strings.Builder
	k := len(alphabet)

	for range n {
		c := alphabet[randGenerator.Intn(k)]
		strBuilder.WriteByte(c)
	}

	return strBuilder.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoeny() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "AUD", "JPY", "GBP", "EGP"}
	n := len(currencies)
	return currencies[randGenerator.Intn(n)]
}