package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/RodBarenco/colab-project-api/db"
)

var forbiddenKeywords = []string{"select", "insert", "update", "delete"}

type validationForCourseAndInstitution struct {
	IsValid    bool
	ExistsInDB bool
}

func IsValidCcourse(ccourse string, currentlyID interface{}, accessor *gorm.DB) validationForCourseAndInstitution {
	result := validationForCourseAndInstitution{}

	if ccourse == "" {
		result.IsValid = true
		return result
	}

	// Check if the course area is a valid string and doesn't contain forbidden symbols.
	if !IsValidField(ccourse) || containsForbiddenSymbol(ccourse) {
		return result
	}

	// Check the length of the course area (you can adjust the minimum and maximum length as needed).
	if len(ccourse) < 3 || len(ccourse) > 50 {
		return result
	}

	// Check if currentlyID is a valid UUID.
	if currentlyIDStr, ok := currentlyID.(string); ok {
		uuidRegex := `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
		if !regexp.MustCompile(uuidRegex).MatchString(currentlyIDStr) {
			return result
		}

		// Check if the institution with the provided ID exists in the database.
		result.ExistsInDB = db.IsInstitutionExists(accessor, uuid.MustParse(currentlyIDStr))
		result.IsValid = true
	}

	return result
}

func IsValidLcourse(lcourse string, lastEducationID interface{}, accessor *gorm.DB) validationForCourseAndInstitution {
	result := validationForCourseAndInstitution{}

	if lcourse == "" {
		result.IsValid = true
		return result
	}

	// Check if the course area is a valid string and doesn't contain forbidden symbols.
	if !IsValidField(lcourse) || containsForbiddenSymbol(lcourse) {
		return result
	}

	// Check the length of the course area (you can adjust the minimum and maximum length as needed).
	if len(lcourse) < 3 || len(lcourse) > 50 {
		return result
	}

	// Check if lastEducationID is a valid UUID.
	if lastEducationIDStr, ok := lastEducationID.(string); ok {
		uuidRegex := `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
		if !regexp.MustCompile(uuidRegex).MatchString(lastEducationIDStr) {
			return result
		}

		// Check if the institution with the provided ID exists in the database.
		result.ExistsInDB = db.IsInstitutionExists(accessor, uuid.MustParse(lastEducationIDStr))
		result.IsValid = true
	}

	return result
}

func IsValidEmail(email string) bool {

	if len(email) > 100 {
		return false
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func IsValidFirstName(firstName string) bool {
	// Check if the first name contains any forbidden keyword
	if containsForbiddenKeyword(firstName) {
		return false
	}

	// Check if the first name matches the required pattern (2 to 25 alphabetical characters)
	firstNameRegex := `^[A-Za-z]{2,25}$`
	return regexp.MustCompile(firstNameRegex).MatchString(firstName)
}

func IsValidLastName(lastName string) bool {
	// Check if the last name contains any forbidden keyword
	if containsForbiddenKeyword(lastName) {
		return false
	}

	// Check if the last name matches the required pattern (1 to 40 alphabetical characters)
	lastNameRegex := `^[A-Za-z]{1,40}$`
	return regexp.MustCompile(lastNameRegex).MatchString(lastName)
}

func IsValidPassword(password string) bool {
	// Check if the password contains any forbidden keyword
	if containsForbiddenKeyword(password) {
		return false
	}

	// Check if the password matches the required pattern (at least 5 characters and no forbidden characters)
	passwordRegex := `^[^\p{C}\/.;:* ]{5,30}$`
	if !regexp.MustCompile(passwordRegex).MatchString(password) {
		return false
	}

	// Check if the password does not contain strings like \\n, \\t or \\r
	if strings.Contains(password, "\\n") || strings.Contains(password, "\\t") || strings.Contains(password, "\\r") {
		return false
	}

	// If passed all checks, return true
	return true
}

func IsValidNickname(nickname string) bool {
	// Check if the nickname contains any forbidden keyword
	if containsForbiddenKeyword(nickname) {
		return false
	}

	// Check if the nickname matches the required pattern (2 to 30 alphanumeric characters, hyphens, or underscores)
	nicknameRegex := `^[A-Za-z0-9_\-]{2,30}$`
	return regexp.MustCompile(nicknameRegex).MatchString(nickname)
}

func IsValidField(field string) bool {
	// Check if the field contains any forbidden symbol
	if containsForbiddenSymbol(field) {
		return false
	}

	if containsForbiddenKeyword(field) {
		return false
	}

	// Check if the field matches the required pattern (alphanumeric characters, hyphens, underscores, commas, periods, spaces, up to 50 characters)
	fieldRegex := `^[A-Za-z0-9_\-,. ]{2,50}$`
	return regexp.MustCompile(fieldRegex).MatchString(field)
}

func IsValidBiography(biography string) bool {
	// Check if the biography contains any forbidden symbol
	if containsForbiddenSymbol(biography) {
		return false
	}

	// Check if the biography matches the required pattern (alphanumeric characters, hyphens, underscores, commas, periods, or spaces, up to 500 characters)
	biographyRegex := `^[A-Za-z0-9_\-,\s.]{3,500}$`
	if !regexp.MustCompile(biographyRegex).MatchString(biography) {
		return false
	}

	// Check if the biography contains any forbidden keywords
	if containsForbiddenKeyword(biography) {
		return false
	}

	return true
}

func IsValidDateOfBirth(dateOfBirth string) bool {
	dateRegex := `^\d{4}-\d{2}-\d{2}$`
	return regexp.MustCompile(dateRegex).MatchString(dateOfBirth)
}

func IsValidInterests(interests []*db.Interest, accessor *gorm.DB) (bool, error) {
	if len(interests) == 0 {
		// No interests provided, valid case
		return true, nil
	}

	for _, interest := range interests {
		// Check if the interest name is a valid string and does not contain forbidden symbols
		if !IsValidField(interest.Name) {
			return false, errors.New("Invalid interest name format")
		}

		// Check if the interest exists in the database
		existsInDB := db.IsInterestExists(accessor, interest.Name)
		if !existsInDB {
			return false, errors.New("One or more interests do not exist or have not been registered")
		}
	}

	return true, nil
}

// ArticleSearch regex

func ArticleSearchIsValid(search string) bool {
	// Check if the search parameter matches the required pattern
	searchRegex := `^[A-Za-z0-9\s]{3,80}$`
	return regexp.MustCompile(searchRegex).MatchString(search)
}

// Função auxiliar para verificar se o campo contém alguma palavra-chave proibida
func containsForbiddenKeyword(field string) bool {
	lowercaseField := strings.ToLower(field)
	for _, keyword := range forbiddenKeywords {
		if strings.Contains(lowercaseField, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// Função auxiliar para verificar se o campo contém algum símbolo proibido
func containsForbiddenSymbol(field string) bool {
	for _, symbol := range []string{"*", ";", "="} {
		if strings.Contains(field, symbol) {
			return true
		}
	}
	return false
}
