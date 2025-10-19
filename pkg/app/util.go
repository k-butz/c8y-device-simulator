package app

import "math/rand/v2"

var direction = [2]int{-1, 1}

func RandFloat(min, max, jitterPercent float64) float64 {
	res := min + rand.Float64()*(max-min)
	if jitterPercent > 0 {
		jitter := res * jitterPercent / 100.0
		direction := direction[rand.IntN(2)]
		res = res + float64(direction)*jitter
	}
	return res
}

func RandInt(min, max int, jitterPercent float64) int {
	res := rand.IntN(max-min+1) + min
	if jitterPercent > 0 {
		jitter := float64(res) * jitterPercent / 100.0
		direction := direction[rand.IntN(2)]
		res = res + direction*int(jitter)
	}
	return res
}
