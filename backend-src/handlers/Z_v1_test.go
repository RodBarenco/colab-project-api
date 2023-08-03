// I cant call handlers_test case I need to import a variable that is just in handlers package and cant be exported
package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// -----------------------PREPARE REFERENCES-------------------------------------------

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
	if status := rr.Code; status != expectedStatus {
		t.Errorf("\n\033[35mExpected status code\033[0m %v \033[35mbut got:\033[0m %v \n", expectedStatus, status)
		return false
	}
	return true
}

func assertResponseBody(t *testing.T, rr *httptest.ResponseRecorder, expectedResponse string) bool {
	t.Helper()
	if rr.Body.String() != expectedResponse {
		t.Errorf("\n\033[35m>>>Expected response body:\033[0m\n %v \n\033[35m>>> BUT GOT:\033[0m\n %v \n", expectedResponse, rr.Body.String())
		return false
	}
	return true
}

// -----------------------TEST FUNCTIONS -------------------------------------------
func TestRegisterHandler_Standard(t *testing.T) {

	cleanup := setupTestEnvironment(t)
	defer cleanup()

	requestBody := StandardRequestBody

	req, err := http.NewRequest("POST", "localhost:8080/v1/register", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to capture the handler's response
	rr := httptest.NewRecorder()

	// CALL TESTED FUNCTION !!!
	RegisterHandler(rr, req)

	// Verify the responses...
	codeRes := assertResponseStatusCode(t, rr, http.StatusCreated)

	expectedResponse := `{"user":{"first_name":"Joe","last_name":"Doe","email":"exemple@example.com"},"message":"User registered!"}`
	bodyRes := assertResponseBody(t, rr, expectedResponse)

	t.Log("\n\nTEST 1 - Standard User Registration !!!\n")
	printTestResults(t, codeRes, bodyRes)
}

func TestRegisterHandler_DuplicateUser(t *testing.T) {

	cleanup := setupTestEnvironment(t)
	defer cleanup()

	requestBody := StandardRequestBody

	req, err := http.NewRequest("POST", "localhost:8000/v1/register", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// CALL TESTED FUNCTION !!!
	RegisterHandler(rr, req)

	// Verify the responses....
	codeRes := assertResponseStatusCode(t, rr, http.StatusInternalServerError)

	expectedResponse := `{"error":"Internal Server Error: Error during signup: failed to create the user: ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)"}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %v but got %v", expectedResponse, rr.Body.String())
	}

	bodyRes := assertResponseBody(t, rr, expectedResponse)

	t.Log("\n\nTEST 2 - Duplicate User Registration !!!\n")
	printTestResults(t, codeRes, bodyRes)
}
