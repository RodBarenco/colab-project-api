package db

import (
	"errors"
	"strings"
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

func ApproveAdmin(db *gorm.DB, adminIDToApprove uuid.UUID, approvingAdminID uuid.UUID) error {
	// Get the admin to approve
	adminToApprove := Admin{}
	err := db.Where("ID = ?", adminIDToApprove).First(&adminToApprove).Error
	if err != nil {
		return err
	}

	// Get the approving admin
	approvingAdmin := Admin{}
	err = db.Where("ID = ?", approvingAdminID).First(&approvingAdmin).Error
	if err != nil {
		return err
	}

	// Check if the approving admin has permission 3 or 4
	if approvingAdmin.Permissions != 3 && approvingAdmin.Permissions != 4 {
		return errors.New("You do not have permission to approve new administrators")
	}

	// Set IsAccepted to true for the admin to approve
	adminToApprove.IsAccepted = true

	// Save the changes to the database
	err = db.Save(&adminToApprove).Error
	if err != nil {
		return err
	}

	return nil
}

func DisapproveAdmin(db *gorm.DB, adminIDToDisapprove uuid.UUID, disapprovingAdminID uuid.UUID) error {
	// Get the admin to disapprove
	adminToDisapprove := Admin{}
	err := db.Where("ID = ?", adminIDToDisapprove).First(&adminToDisapprove).Error
	if err != nil {
		return err
	}

	// Get the disapproving admin
	disapprovingAdmin := Admin{}
	err = db.Where("ID = ?", disapprovingAdminID).First(&disapprovingAdmin).Error
	if err != nil {
		return err
	}

	// Check if the disapproving admin has permission 3 or 4
	if disapprovingAdmin.Permissions != 3 && disapprovingAdmin.Permissions != 4 {
		return errors.New("You do not have permission to disapprove administrators")
	}

	// Check if the disapproving admin has a lower permission level than the admin to disapprove
	if disapprovingAdmin.Permissions >= adminToDisapprove.Permissions {
		return errors.New("You cannot disapprove an administrator with equal or higher permissions")
	}

	// Set IsAccepted to false for the admin to disapprove
	adminToDisapprove.IsAccepted = false

	// Save the changes to the database
	err = db.Save(&adminToDisapprove).Error
	if err != nil {
		return err
	}

	return nil
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

func DeleteUser(db *gorm.DB, userIDToDelete uuid.UUID, deletingAdminID uuid.UUID) error {
	// Get the user to delete
	userToDelete := User{}
	err := db.Where("ID = ?", userIDToDelete).First(&userToDelete).Error
	if err != nil {
		return err
	}

	// Get the deleting admin
	deletingAdmin := Admin{}
	err = db.Where("ID = ?", deletingAdminID).First(&deletingAdmin).Error
	if err != nil {
		return err
	}

	// Check if the deleting admin has permission 3 or 4
	if deletingAdmin.Permissions != 3 && deletingAdmin.Permissions != 4 {
		return errors.New("You do not have permission to delete users")
	}

	// Perform the user deletion
	err = db.Delete(&userToDelete).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteAdmin(db *gorm.DB, adminIDToDelete uuid.UUID, deletingAdminID uuid.UUID) error {
	// Get the admin to delete
	adminToDelete := Admin{}
	err := db.Where("ID = ?", adminIDToDelete).First(&adminToDelete).Error
	if err != nil {
		return err
	}

	// Get the deleting admin
	deletingAdmin := Admin{}
	err = db.Where("ID = ?", deletingAdminID).First(&deletingAdmin).Error
	if err != nil {
		return err
	}

	// Check if the deleting admin has permission 3 or 4
	if deletingAdmin.Permissions != 3 && deletingAdmin.Permissions != 4 {
		return errors.New("You do not have permission to delete administrators")
	}

	// Check if the admin to delete has a lower permission level than the deleting admin
	if adminToDelete.Permissions >= deletingAdmin.Permissions {
		return errors.New("You cannot delete an administrator with equal or higher permissions")
	}

	// Check if the admin to delete is not accepted
	if adminToDelete.IsAccepted {
		return errors.New("You can only delete administrators who are not accepted")
	}

	// Perform the admin deletion
	err = db.Delete(&adminToDelete).Error
	if err != nil {
		return err
	}

	return nil
}

func GetUnacceptedArticlesID(db *gorm.DB, dateToCompare time.Time) ([]uuid.UUID, error) {
	// Query for unaccepted articles based on the provided date
	var unacceptedArticleIDs []uuid.UUID
	err := db.Table("articles").
		Select("id").
		Where("is_accepted = ? AND acceptance_date <= ?", false, dateToCompare).
		Pluck("id", &unacceptedArticleIDs).
		Error
	if err != nil {
		return nil, err
	}

	return unacceptedArticleIDs, nil
}

func GetUnacceptedArticlesIDByField(db *gorm.DB, fieldToCompare string) ([]uuid.UUID, error) {
	// Query for unaccepted articles based on the provided field (case-insensitive)
	var unacceptedArticleIDs []uuid.UUID
	err := db.Table("articles").
		Select("id").
		Where("is_accepted = ? AND LOWER(field) = LOWER(?)", false, fieldToCompare).
		Pluck("id", &unacceptedArticleIDs).
		Error
	if err != nil {
		return nil, err
	}

	return unacceptedArticleIDs, nil
}

// just root
func CleanOldUnacceptedArticlesByDate(db *gorm.DB, dateToCompare time.Time, rootAdminID uuid.UUID) error {
	// Get the root admin
	rootAdmin := Admin{}
	err := db.Where("ID = ?", rootAdminID).First(&rootAdmin).Error
	if err != nil {
		return err
	}

	// Check if the root admin has permission 4
	if rootAdmin.Permissions != 4 {
		return errors.New("Only root administrators can clean old unaccepted articles")
	}

	// Clean old unaccepted articles based on the provided date
	err = db.Where("is_accepted = ? AND acceptance_date <= ?", false, dateToCompare).
		Delete(&Article{}).Error

	if err != nil {
		return err
	}

	return nil
}

func CleanOldUnacceptedArticlesByDateAndField(db *gorm.DB, dateToCompare time.Time, field string, adminID uuid.UUID) error {
	// Get the admin
	admin := Admin{}
	err := db.Where("ID = ?", adminID).First(&admin).Error
	if err != nil {
		return err
	}

	// Convert 'field' to lowercase
	field = strings.ToLower(field)

	// Check if the admin has permission 3 and 'field' matches
	if admin.Permissions == 3 && strings.ToLower(admin.Field) == field {
		// Clean old unaccepted articles based on the provided date and field
		err = db.Where("is_accepted = ? AND acceptance_date <= ? AND lower(field) = ?",
			false, dateToCompare, field).
			Delete(&Article{}).Error

		if err != nil {
			return err
		}

		return nil
	}

	// Check if the admin has permission 4 (root)
	if admin.Permissions == 4 {
		// Clean old unaccepted articles based on the provided date and field
		err = db.Where("is_accepted = ? AND acceptance_date <= ? AND lower(field) = ?",
			false, dateToCompare, field).
			Delete(&Article{}).Error

		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("You do not have permission to clean old unaccepted articles with this field")
}

// just root
func CleanAllOldUnacceptedArticles(db *gorm.DB, rootAdminID uuid.UUID) error {
	// Get the root admin
	rootAdmin := Admin{}
	err := db.Where("ID = ?", rootAdminID).First(&rootAdmin).Error
	if err != nil {
		return err
	}

	// Check if the root admin has permission 4
	if rootAdmin.Permissions != 4 {
		return errors.New("Only root administrators can clean all old unaccepted articles")
	}

	// Clean all old unaccepted articles
	err = db.Where("is_accepted = ?", false).
		Delete(&Article{}).Error

	if err != nil {
		return err
	}

	return nil
}

// GET PKEY

func GetAdminPublicKey(db *gorm.DB, AdminID uuid.UUID) (string, error) {
	var admin Admin
	result := db.Select("public_key").Where("id = ?", AdminID).First(&admin)
	if result.Error != nil {
		return "", result.Error
	}
	return admin.PublicKey, nil
}
