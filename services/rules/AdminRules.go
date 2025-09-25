package rules

import (
    "sort"

    "github.com/apkatsikas/artist-entities/storageclient"
)

type AdminRules struct {
}

const max = 2

func (ar *AdminRules) FileToDelete(files []storageclient.BackupFile) string {
    if len(files) > max {
        // Sort by updated date, oldest first
        sort.Slice(files[:], func(i, j int) bool {
            return files[i].Updated.Before(files[j].Updated)
        })

        // Return oldest entry's name
        return files[0].Name
    }
    // Return blank if less files than the max
    // Nothing to delete
    return ""
}
