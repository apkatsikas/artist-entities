package rules

import (
    "math/rand"
    "regexp"
    "strings"
    "time"

    ce "github.com/apkatsikas/artist-entities/customerrors"
)

type ArtistRules struct {
}

const (
    limit        = 75
    theWithSpace = "the "
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func (rules *ArtistRules) CleanArtistName(artistName string) (string, error) {
    badChars := []string{"\"", ","}

    // Return an error if any bad characters
    for _, bad := range badChars {
        if strings.Contains(artistName, bad) {
            return "", ce.ErrDataInvalid
        }
    }
    // If we have more than one single quote (1 is fine)
    if strings.Count(artistName, "'") > 1 {
        return "", ce.ErrDataInvalid
    }

    // Trim space
    artistName = strings.TrimSpace(artistName)
    // Lowercase
    artistName = strings.ToLower(artistName)
    // Remove first the
    if strings.HasPrefix(artistName, theWithSpace) {
        artistName = strings.Replace(artistName, theWithSpace, "", 1)
    }

    // Remove all special characters
    artistName = nonAlphanumericRegex.ReplaceAllString(artistName, "")

    // Check length AFTER we manipulate
    // as we may truncate to the limit during clean
    if len(artistName) > limit {
        return "", ce.ErrDataTooLong
    }
    // If we removed all characters, return invalid
    if len(artistName) <= 0 {
        return "", ce.ErrDataInvalid
    }
    return artistName, nil
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
