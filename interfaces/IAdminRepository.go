package interfaces

type IAdminRepository interface {
    CreateBackup(file string) error
}
