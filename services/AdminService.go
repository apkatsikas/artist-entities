package services

import (
    "fmt"
    "time"

    "github.com/apkatsikas/artist-entities/interfaces"
)

const (
    vacuumFileName = "vacuum.sqlite"
    entitiesBackup = "entities-backup"
)

type AdminService struct {
    AdminRepository interfaces.IAdminRepository
    FileUtil        interfaces.IFileUtil
    StorageClient   interfaces.IStorageClient
    Rules           interfaces.IAdminRules
}

func (as *AdminService) Backup() error {
    // Remove existing vacuum file first
    err := as.FileUtil.DeleteIfExists(vacuumFileName)
    if err != nil {
        return err
    }

    // Back up the DB
    err = as.AdminRepository.CreateBackup(vacuumFileName)
    if err != nil {
        return err
    }

    // Get timestamped file name
    timestamp := time.Now().Unix()
    fileName := fmt.Sprintf("%v%v.sqlite", entitiesBackup, timestamp)

    // Upload file
    err = as.StorageClient.UploadFile(vacuumFileName, fileName)
    if err != nil {
        return err
    }

    // List files
    files, err := as.StorageClient.ListFiles()
    if err != nil {
        return err
    }

    // Determine if we need to delete
    fileToDelete := as.Rules.FileToDelete(files)
    if fileToDelete != "" {
        // Delete
        err = as.StorageClient.DeleteFile(fileToDelete)
        if err != nil {
            return err
        }
    }

    return nil
}
