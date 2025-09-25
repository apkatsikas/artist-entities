package interfaces

import "github.com/apkatsikas/artist-entities/storageclient"

type IStorageClient interface {
    DeleteFile(object string) error
    UploadFile(path string, destObject string) error
    ListFiles() ([]storageclient.BackupFile, error)
}
