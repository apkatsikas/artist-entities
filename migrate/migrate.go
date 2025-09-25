package migrate

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/apkatsikas/artist-entities/customerrors"
	"github.com/apkatsikas/artist-entities/infrastructures/logutil"
	"github.com/apkatsikas/artist-entities/interfaces"
)

const filePath = "artists.csv"

func readCsvFile() [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		logutil.Error("Unable to read input file, error was %v", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		logutil.Error("Unable to parse input file, error was %v", err)
	}

	return records
}

func Migrate(ar interfaces.IArtistRepository, as interfaces.IArtistService, secret string) {
	// Setup table
	err := ar.Migrate()
	if err != nil {
		logutil.Error("Got an unexpected error during artist migration: %v", err)
	}

	artists := readCsvFile()

	// Insert records
	for _, a := range artists[0] {
		result, err := as.Create(a)
		if err != nil {
			// Log if record exists, or something unexpected
			if errors.Is(err, customerrors.ErrRecordExists) {
				log := fmt.Sprintf("Artist '%v' already exists", a)
				logutil.Info(log)
			} else {
				logutil.Error("Got an unexpected error on record %v, error was: %v", a, err)
			}

		} else {
			log := fmt.Sprintf("Inserted artist: %v", result)
			logutil.Info(log)
		}
	}
}
