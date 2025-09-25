package models

import "gorm.io/gorm"

type Artist struct {
    gorm.Model
    Name string `gorm:"type:varchar(75);unique_index;not null"`
}
