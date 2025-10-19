package app

import (
	"math/rand/v2"
	"time"
)

const TS_FORMAT string = "2006-01-02T15:04:05.000Z07:00"

var direction = [2]int{-1, 1}

func RandFloat(min, max, jitterPercent float64) float64 {
	res := min + rand.Float64()*(max-min)
	if jitterPercent > 0 {
		jitterMin := 0.0
		jitterMax := res * jitterPercent / 100.0
		rndJitter := RandFloat(jitterMin, jitterMax, 0)
		direction := direction[rand.IntN(2)]
		res = res + float64(direction)*rndJitter
	}
	return res
}

func ToRFCTimeStamp(t time.Time) string {
	return t.Format(TS_FORMAT)
}
