// I cant call handlers_test case I need to import a variable that is just in handlers package and cant be exported
package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RodBarenco/colab-project-api/utils"
)

// -----------------------TEST FUNCTIONS FOR V1-------------------------------------------
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

	t.Log(utils.OrangeColor.InitColor + "\n\nTEST 1 - Standard User Registration !!!\n" + utils.EndColor)
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

	t.Log(utils.OrangeColor.InitColor + "\n\nTEST 2 - Duplicate User Registration !!!\n" + utils.EndColor)
	printTestResults(t, codeRes, bodyRes)
}

// TEST FIRST NAME
func TestRegisterHandlerWithInvalidFirstName(t *testing.T) {
	invalidFirstNames := []string{
		"",
		"select",
		"Joe*",
		"J",
		"JoeDoeJoeDoeJoeDoeJoeDoeJoeDoe",
	}

	cleanup := setupTestEnvironment(t)
	defer cleanup()

	for i, firstName := range invalidFirstNames {
		t.Run(fmt.Sprintf("InvalidFirstNameCase_%d", i+1), func(t *testing.T) {
			// Generate a RequestBody with FirstName set to the test value
			requestBody := GenerateRequestBodyWithFirstName(firstName)

			req, err := http.NewRequest("POST", "localhost:8080/v1/register", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a ResponseRecorder to capture the handler's response
			rr := httptest.NewRecorder()

			// CALL TESTED FUNCTION !!!
			RegisterHandler(rr, req)

			// Verify the responses...
			codeRes := assertResponseStatusCode(t, rr, http.StatusBadRequest)

			expectedResponse := `{"error":"Bad Request: {1 : First name must have 2 to 25 characters - and valid characters-words}"}`
			bodyRes := assertResponseBody(t, rr, expectedResponse)

			t.Logf(utils.OrangeColor.InitColor+"\n\nTEST INVALID FIRST NAME %d - User Registration with invalid first name format: %s !!!\n"+utils.EndColor, i+1, firstName)
			printTestResults(t, codeRes, bodyRes)
		})
	}
}

// TEST LAST NAME
func TestRegisterHandlerWithInvalidLastName(t *testing.T) {
	invalidLastNames := []string{
		"",
		"select",
		"Doe*",
		"Doeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	}

	cleanup := setupTestEnvironment(t)
	defer cleanup()

	for i, lastName := range invalidLastNames {
		t.Run(fmt.Sprintf("InvalidLastNameCase_%d", i+1), func(t *testing.T) {
			// Generate a RequestBody with LastName set to the test value
			requestBody := GenerateRequestBodyWithLastName(lastName)

			req, err := http.NewRequest("POST", "localhost:8080/v1/register", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a ResponseRecorder to capture the handler's response
			rr := httptest.NewRecorder()

			// CALL TESTED FUNCTION !!!
			RegisterHandler(rr, req)

			// Verify the responses...
			codeRes := assertResponseStatusCode(t, rr, http.StatusBadRequest)

			expectedResponse := `{"error":"Bad Request: {1 : Last name must have 1 to 40 characters - and valid characters-words}"}`
			bodyRes := assertResponseBody(t, rr, expectedResponse)

			t.Logf(utils.OrangeColor.InitColor+"\n\nTEST INVALID LAST NAME %d - User Registration with invalid last name format: %s !!!\n"+utils.EndColor, i+1, lastName)
			printTestResults(t, codeRes, bodyRes)
		})
	}
}

// TEST EMAIL

func TestRegisterHandlerWithValidSpecialCharInMail(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Generate a RequestBody with FirstName set to an empty string (null)
	requestBody := GenerateRequestBodyWithEmial("joe_doe@example.com")

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

	expectedResponse := `{"user":{"first_name":"Joe","last_name":"Doe","email":"joe_doe@example.com"},"message":"User registered!"}`
	bodyRes := assertResponseBody(t, rr, expectedResponse)

	t.Log(utils.OrangeColor.InitColor + "\n\nTEST VALID EMAIL 1 - User Registration with valid special character in emaill !!!\n" + utils.EndColor)
	printTestResults(t, codeRes, bodyRes)
}

func TestRegisterHandlerWithInvalidEmail(t *testing.T) {
	invalidEmails := []string{
		"",
		"joedoeexemple.com",
		"joedoe@exemplecom",
		"joedoe@example.c",
		"joedoe!@example.com",
		"joedoe@example!.com",
		"joe doe@example.com",
		"joedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoejoedoe@exemple.com",
	}

	cleanup := setupTestEnvironment(t)
	defer cleanup()

	for i, email := range invalidEmails {
		t.Run(fmt.Sprintf("InvalidEmailCase_%d", i+1), func(t *testing.T) {
			// Generate a RequestBody with FirstName set to an empty string (null)
			requestBody := GenerateRequestBodyWithEmial(email)

			req, err := http.NewRequest("POST", "localhost:8080/v1/register", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a ResponseRecorder to capture the handler's response
			rr := httptest.NewRecorder()

			// CALL TESTED FUNCTION !!!
			RegisterHandler(rr, req)

			// Verify the responses...
			codeRes := assertResponseStatusCode(t, rr, http.StatusBadRequest)

			expectedResponse := `{"error":"Bad Request: {1 : Invalid email format}"}`
			bodyRes := assertResponseBody(t, rr, expectedResponse)

			t.Logf(utils.OrangeColor.InitColor+"\n\nTEST INVALID EMAIL %d - User Registration with invalid email format: %s !!!\n"+utils.EndColor, i+1, email)
			printTestResults(t, codeRes, bodyRes)
		})
	}
}

// TEST PASSWORD
func TestRegisterHandlerWithInvalidPassword(t *testing.T) {
	invalidPasswords := []string{
		"Ab@1",
		"ThisIsAnExtremelyLongPassword1234567890!@#",
		"password; DROP TABLE users;--",
		"ThisPassContainsSpaces ",
		" MyPassHasLeadingSpaces",
		"MyPassHasInvalidSymbols\000\001",
		"MyPassHasInvalidSymbols\\n\\t\\r",
		"ThisPassIsTooLongItExceedsTheMaxAllowedLimitOfThirtyCharacters",
	}

	cleanup := setupTestEnvironment(t)
	defer cleanup()

	for i, password := range invalidPasswords {
		t.Run(fmt.Sprintf("InvalidPasswordCase_%d", i+1), func(t *testing.T) {
			// Generate a RequestBody with Password set to the test value
			requestBody := GenerateRequestBodyWithPassWord(password)

			req, err := http.NewRequest("POST", "localhost:8080/v1/register", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create a ResponseRecorder to capture the handler's response
			rr := httptest.NewRecorder()

			// CALL TESTED FUNCTION !!!
			RegisterHandler(rr, req)

			// Verify the responses...
			codeRes := assertResponseStatusCode(t, rr, http.StatusBadRequest)

			expectedResponse := `{"error":"Bad Request: {1 : Password must have at least 5 characters - and valid characters-words}"}`
			bodyRes := assertResponseBody(t, rr, expectedResponse)

			t.Logf(utils.OrangeColor.InitColor+"\n\nTEST INVALID PASSWORD %d - User Registration with invalid password format: %s !!!\n"+utils.EndColor, i+1, password)
			printTestResults(t, codeRes, bodyRes)
		})
	}
}
