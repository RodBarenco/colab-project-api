package db

import "gorm.io/gorm"

type Image struct {
	ID   uint   `gorm:"primaryKey"`
	Data string `gorm:"not null"`
}

func SaveImageToDB(db *gorm.DB, imageBase64 string) (uint, error) {
	newImage := Image{
		Data: imageBase64,
	}

	if err := db.Create(&newImage).Error; err != nil {
		return 0, err
	}

	return newImage.ID, nil
}

func GetImageBase64ByID(db *gorm.DB, imageID uint) (string, error) {
	var image Image
	if err := db.First(&image, imageID).Error; err != nil {
		return "", err
	}

	return image.Data, nil
}
