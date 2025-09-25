package rules

import (
    "testing"
    "time"

    sc "github.com/apkatsikas/artist-entities/storageclient"
    "github.com/stretchr/testify/assert"
)

func TestAdminRules(t *testing.T) {
    var testData = []struct {
        files        []sc.BackupFile
        fileToDelete string
        test         string
    }{
        {
            test:         "no files - nothing to delete",
            files:        []sc.BackupFile{},
            fileToDelete: "",
        },
        {
            test:         "1 file - nothing to delete",
            files:        []sc.BackupFile{{Name: "a", Updated: time.Now()}},
            fileToDelete: "",
        },
        {
            test: "2 files - nothing to delete",
            files: []sc.BackupFile{
                {Name: "a", Updated: time.Now()},
                {Name: "b", Updated: time.Now()},
            },
            fileToDelete: "",
        },
        {
            test: "3 files - delete oldest - middle",
            files: []sc.BackupFile{
                {Name: "a", Updated: time.Now()},
                {Name: "b", Updated: time.Now().AddDate(0, -2, 0)},
                {Name: "c", Updated: time.Now().AddDate(0, -1, 0)},
            },
            fileToDelete: "b",
        },
        {
            test: "3 files - delete oldest - first",
            files: []sc.BackupFile{
                {Name: "a", Updated: time.Now().AddDate(0, -2, 0)},
                {Name: "b", Updated: time.Now()},
                {Name: "c", Updated: time.Now().AddDate(0, -1, 0)},
            },
            fileToDelete: "a",
        },
        {
            test: "3 files - delete oldest - last",
            files: []sc.BackupFile{
                {Name: "a", Updated: time.Now().AddDate(0, -1, 0)},
                {Name: "b", Updated: time.Now()},
                {Name: "c", Updated: time.Now().AddDate(0, -2, 0)},
            },
            fileToDelete: "c",
        },
        {
            test: "4 files - delete oldest - last",
            files: []sc.BackupFile{
                {Name: "a", Updated: time.Now().AddDate(0, -1, 0)},
                {Name: "b", Updated: time.Now()},
                {Name: "c", Updated: time.Now().AddDate(0, -2, 0)},
                {Name: "d", Updated: time.Now().AddDate(0, -3, 0)},
            },
            fileToDelete: "d",
        },
    }
    for _, tt := range testData {
        t.Run(tt.test, func(t *testing.T) {
            rules := AdminRules{}
            result := rules.FileToDelete(tt.files)
            assert.Equal(t, result, tt.fileToDelete)
        })
    }
}
