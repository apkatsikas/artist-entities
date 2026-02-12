package rules

import (
	"math/rand"
	"time"
)

type ArtistRules struct {
}

func (rules *ArtistRules) RandomOffset(count uint) uint {
	min := 1
	countInt := int(count)

	if countInt == min {
		return uint(min)
	}

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	maxMinusMin := countInt - min

	offset := r.Intn(maxMinusMin) + min
	return uint(offset)
}
