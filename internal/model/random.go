package model

type RandomSupport struct {
	Float64 func() float64 // Float64 returns a random float64 in [0.0,1.0).
}

// RandomlyDecide returns true if the random number is less than howLikely. For example, randomlyDecide(0.1) will return true 10% of the time.
func (rs RandomSupport) RandomlyDecide(howLikely float64) bool {
	if howLikely < 0 || howLikely > 1 {
		panic("howLikely must be in the range 0.0 to 1.0")
	}
	return rs.Float64() < howLikely
}

// randomlyChange returns a value that is randomly changed by a certain percentage. For example, randomlyChange(100, 0.1) will return a value between 90 and 110.
// randomlyChange can increase or decrease the value.
func (rs RandomSupport) RandomlyChange(val float64, howMuch float64) float64 {
	if howMuch < 0 || howMuch > 1 {
		panic("howMuch must be in the range 0.0 to 1.0")
	}

	rndFactor := 1 - rs.Float64()*2 // Random number in half-open interval (-1, 1]
	return val * (1 + rndFactor*howMuch)
}
