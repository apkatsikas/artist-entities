package repositories

import (
    "fmt"

    "github.com/apkatsikas/artist-entities/interfaces"
)

type AdminRepository struct {
    IDB interfaces.IDbHandler
}

func (adR *AdminRepository) CreateBackup(file string) error {
    gormConn := adR.IDB.Connection()

    vQ := fmt.Sprintf("VACUUM main into '%v';", file)
    vRes := gormConn.Exec(vQ)

    if vRes.Error != nil {
        return vRes.Error
    }
    return nil
}
