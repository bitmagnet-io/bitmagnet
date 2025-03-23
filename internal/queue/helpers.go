package queue

import (
	"math"
	"math/rand"
	"time"
)

// RandInt returns a random integer up to max
func RandInt(maxValue int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(maxValue)
}

// CalculateBackoff calculates the number of seconds to back off before the next retry
// this formula is unabashedly taken from Sidekiq because it is good.
func CalculateBackoff(retryCount uint) time.Time {
	const backoffExponent = 4

	const maxInt = 30

	p := int(math.Round(math.Pow(float64(retryCount), backoffExponent)))

	return time.Now().UTC().Add(time.Duration(p+15+RandInt(maxInt)*int(retryCount)+1) * time.Second)
}
