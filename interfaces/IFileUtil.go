package interfaces

type IFileUtil interface {
    DeleteIfExists(file string) error
}
