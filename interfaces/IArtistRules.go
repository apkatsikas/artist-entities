package interfaces

type IArtistRules interface {
    CleanArtistName(s string) (string, error)
    RandomOffset(count uint) uint
}
