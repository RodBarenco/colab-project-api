package db

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Institution{}, &Article{}, &Image{})
}
