package interfaces

import "gorm.io/gorm"

type IDbHandler interface {
    Connection() *gorm.DB
}
