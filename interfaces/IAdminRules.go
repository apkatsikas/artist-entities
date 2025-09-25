package interfaces

import "github.com/apkatsikas/artist-entities/storageclient"

type IAdminRules interface {
    FileToDelete(files []storageclient.BackupFile) string
}
