package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name string `gorm:"type:varchar(75);unique_index;not null"`
    Password string `gorm:"type:varchar(75);not null"`
}
