package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permissions are 0 up to 4 - the first Admin get permission 4 it means that he can to every action,
// then 0 is the standard and is the one that have more limitations. it is given to all new admins
type Admin struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName    string    `gorm:"not null"`
	LastName     string    `gorm:"not null"`
	Nickname     string
	Email        string    `gorm:"not null;unique"`
	Password     string    `gorm:"not null"`
	DateOfBirth  time.Time `gorm:"not null"`
	Title        string
	Field        string
	Biography    string      `gorm:"type:TEXT"`
	Currently    Institution `gorm:"foreignKey:CurrentlyID"`
	CurrentlyID  *uuid.UUID  `gorm:"null;type:uuid"`
	CreatedAt    time.Time   `gorm:"autoCreateTime"`
	PublicKey    string      `gorm:"not null"`
	IsAccepted   bool
	Permissions  uint `gorm:"not null"`
	ProfilePhoto string
}

func ApproveArticle(db *gorm.DB, articleID uuid.UUID, adminID uuid.UUID) (string, error) {
	// Get the article and admin
	article := Article{}
	err := db.Where("ID = ?", articleID).First(&article).Error
	if err != nil {
		return "", err
	}

	admin := Admin{}
	err = db.Where("ID = ?", adminID).First(&admin).Error
	if err != nil {
		return "", err
	}

	// Check if the admin is root
	if admin.Permissions == 4 {
		// If the admin is root, approve the article directly
		article.IsAccepted = true
		article.AcceptanceDate = time.Now()
		// Add the root to the list of approved admins
		article.ApprovedBy = append(article.ApprovedBy, &admin)
		err = db.Save(&article).Error
		if err != nil {
			return "", err
		}

		return "Approved by root admin", nil
	} else {
		// If the admin is not root, add them to the list of approved admins
		article.ApprovedBy = append(article.ApprovedBy, &admin)
		err = db.Save(&article).Error
		if err != nil {
			return "", err
		}

		// Check the number of approved admins
		numApprovedAdmins := len(article.ApprovedBy)

		// If the admin is level 0 or 1, we need 3 approvals
		if admin.Permissions <= 1 {
			if numApprovedAdmins >= 3 {
				article.IsAccepted = true
				article.AcceptanceDate = time.Now()
				err = db.Save(&article).Error
				if err != nil {
					return "", err
				}
			}
		} else {
			// If the admin is level 2 or 3, we need 2 approvals
			if numApprovedAdmins >= 2 {
				article.IsAccepted = true
				article.AcceptanceDate = time.Now()
				err = db.Save(&article).Error
				if err != nil {
					return "", err
				}
			}
		}
	}

	return "Admin " + admin.ID.String() + " was added to approvedBy", nil
}

func DeleteArticle(db *gorm.DB, articleID uuid.UUID, adminID uuid.UUID) error {
	// Get the article
	article := Article{}
	err := db.Where("ID = ?", articleID).First(&article).Error
	if err != nil {
		return err
	}

	// Get admin
	admin := Admin{}
	err = db.Where("ID = ?", adminID).First(&admin).Error
	if err != nil {
		return err
	}

	// Check if the admin has permission 0
	if admin.Permissions == 0 {
		return errors.New("You do not have permission to delete this article")
	}

	// Check if the admin is root or has permission 3
	if admin.Permissions == 4 || admin.Permissions == 3 {
		// Delete the article
		err = db.Delete(&article).Error
		if err != nil {
			return err
		}

		return nil
	}

	// Check if the admin has permission 1 or 2
	if admin.Permissions <= 2 {
		// Check if the article is not approved
		if !article.IsAccepted {
			// Delete the article
			err = db.Delete(&article).Error
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("You do not have permission to delete this article")
}

func ApproveAdmin() {
	// TODO
}

func ModifyAdminPermissions(db *gorm.DB, allower uuid.UUID, allowed uuid.UUID, permissions uint) error {
	// Get the current allower.

	admAllower := Admin{}
	err := db.Where("ID = ?", allower).First(&admAllower).Error
	if err != nil {
		return err
	}

	// Get the current allower.

	admAllowed := Admin{}
	err = db.Where("ID = ?", allower).First(&admAllowed).Error
	if err != nil {
		return err
	}

	// Check if the current user has the necessary permissions.
	if admAllower.Permissions < 4 {
		return errors.New("unauthorized")
	}

	// Check if the permissions value is valid.
	if permissions < 0 || permissions > 4 {
		return errors.New("invalid permissions value")
	}

	// Update the admin's permissions.
	admAllowed.Permissions = permissions
	err = db.Save(admAllowed).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser() {
	// TODO
}

func DeleteAdmin() {
	// TODO
}

func GetUnacceptedArticles() {
	// TODO
}

func GetUnacceptedArticlesByField() {
	// TODO
}

// just root
func CleanOldUnacceptedArticlesByDate() {
	// TODO
}

func CleanOldUnacceptedArticlesByDateAndField() {
	// TODO
}

// just root
func CleanAllOldUnacceptedArticles() {
	// TODO
}
