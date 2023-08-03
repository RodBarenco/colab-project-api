// I cant call handlers_test case I need to import a variable that is just in handlers package and cant be exported
package handlers

import (
	"bytes"
	"encoding/json"
	"log"

	"net/http/httptest"
	"os"
	"testing"

	"github.com/RodBarenco/colab-project-api/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// -----------------------PREPARE REFERENCES-------------------------------------------

// ------------BODIES -------
type TestRequestBody struct {
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	Email       string `json:"Email"`
	Password    string `json:"Password"`
	DateOfBirth string `json:"DateOfBirth"`
	Nickname    string `json:"Nickname"`
	Field       string `json:"Field"`
	Biography   string `json:"Biography"`
	OpenToColab bool   `json:"OpenToColab"`
}

var StandardRequestBody = []byte(`{
	"FirstName": "Joe",
	"LastName": "Doe",
	"Email": "exemple@example.com",
	"Password": "123456",
	"DateOfBirth": "1990-01-01",
	"Nickname": "johndoe",
	"Field": "Software Engineering",
	"Biography": "A passionate software engineer.",
	"OpenToColab": true
}`)

// sets the first name in the request body.
func GenerateRequestBodyWithFirstName(firstName string) []byte {
	var body TestRequestBody
	err := json.Unmarshal(StandardRequestBody, &body)
	if err != nil {
		panic("Failed to unmarshal standard request body")
	}
	body.FirstName = firstName
	rb, err := json.Marshal(body)
	if err != nil {
		panic("Failed to marshal modified request body")
	}
	return rb
}

func GenerateRequestBodyWithLastName(lastName string) []byte {
	var body TestRequestBody
	err := json.Unmarshal(StandardRequestBody, &body)
	if err != nil {
		panic("Failed to unmarshal standard request body")
	}
	body.LastName = lastName
	rb, err := json.Marshal(body)
	if err != nil {
		panic("Failed to marshal modified request body")
	}
	return rb
}

func GenerateRequestBodyWithEmial(email string) []byte {
	var body TestRequestBody
	err := json.Unmarshal(StandardRequestBody, &body)
	if err != nil {
		panic("Failed to unmarshal standard request body")
	}
	body.Email = email
	rb, err := json.Marshal(body)
	if err != nil {
		panic("Failed to marshal modified request body")
	}
	return rb
}

func GenerateRequestBodyWithPassWord(password string) []byte {
	var body TestRequestBody
	err := json.Unmarshal(StandardRequestBody, &body)
	if err != nil {
		panic("Failed to unmarshal standard request body")
	}
	body.Password = password
	rb, err := json.Marshal(body)
	if err != nil {
		panic("Failed to marshal modified request body")
	}
	return rb
}

// ----- GENERAL ------

func setupTestEnvironment(t *testing.T) func() {
	err := godotenv.Load("../.test.env")
	if err != nil {
		log.Fatalf("failed to load .test.env")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not found in the environment")
	}
	// discart log messages from DB
	var buf bytes.Buffer
	customLogger := logger.New(
		log.New(&buf, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             0,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// gorm uses cunstonlogger
	dbAccessTest, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		panic("Failed to connect to the database")
	}

	dbAccessor = dbAccessTest
	return func() {
		// Add any cleanup logic here if needed ......
	}
}

func printTestResults(t *testing.T, codeRes, bodyRes bool) {
	if codeRes {
		t.Log("\n \033[32mCode Response has PASSED!!!\033[0m\n")
	} else {
		t.Log("\n \033[33mCode Response has FAILED!!!\033[0m\n")
	}

	if bodyRes {
		t.Log("\n \033[32mBody Response has PASSED!!!\033[0m\n")
	} else {
		t.Log("\n \033[33mBody Response has FAILED!!!\033[0m\n")
	}
}

func assertResponseStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) bool {
	t.Helper()
	c := utils.PurpleColor.InitColor
	e := utils.EndColor
	if status := rr.Code; status != expectedStatus {
		t.Errorf("\n%sExpected status code%s %v %sbut got:%s %v \n", c, e, expectedStatus, c, e, status)
		return false
	}
	return true
}

func assertResponseBody(t *testing.T, rr *httptest.ResponseRecorder, expectedResponse string) bool {
	t.Helper()
	c := utils.PurpleColor.InitColor
	e := utils.EndColor
	if rr.Body.String() != expectedResponse {
		t.Errorf("\n%s>>>Expected response body:%s\n %v \n%s>>> BUT GOT:%s\n %v \n", c, e, expectedResponse, c, e, rr.Body.String())

		return false
	}
	return true
}
