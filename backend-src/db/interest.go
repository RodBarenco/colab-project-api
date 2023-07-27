package db

import "gorm.io/gorm"

//yet to change it ro uint32,  gorm create many-to-may relation with big int id,
type Interest struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null"`
}

// IsInterestExists checks if an interest with the given name exists in the database.
func IsInterestExists(accessor *gorm.DB, interestName string) bool {
	var interest Interest
	result := accessor.Where("name = ?", interestName).First(&interest)
	return result.Error == nil
}
