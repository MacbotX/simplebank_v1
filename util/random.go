package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "asdfghjklqwertyuiopzxcvbnm"

func init() {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))
}

// RandomInt generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Randon string generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency
func RandomCurrency() string {
	currencies := []string{"USD", "NGN", "GHS"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// random email addresss
func RandomEmail() string {
	// random email provider
	email_providers := []string{"gmail", "yahoo", "private"}
	e := len(email_providers)
	provider := email_providers[rand.Intn(e)]

	// generate random string
	ran_str := RandomString(8)

	return fmt.Sprintf("%s@%s.com", ran_str, provider)
}
