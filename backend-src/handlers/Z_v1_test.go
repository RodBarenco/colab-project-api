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
var invalidOtherTests = make([]string, 0)

func TestRegisterHandlerStandard(t *testing.T) {

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

	if !bodyRes || !codeRes {
		invalidOtherTests = append(invalidOtherTests, fmt.Sprintf("Invalid -> TestRegisterHandlerStandard"))
		testResult.OtherTestsCases = append(testResult.OtherTestsCases, invalidOtherTests)
	} else {
		testResult.PassedTests++
	}

	t.Log(utils.OrangeColor.InitColor + "\n\nTEST 1 - Standard User Registration !!!\n" + utils.EndColor)
	printTestResults(t, codeRes, bodyRes)

	testResult.TotalTests++
	testResult.GeneralTests++

}

func TestRegisterHandlerDuplicateUser(t *testing.T) {

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

	if !bodyRes || !codeRes {
		invalidOtherTests = append(invalidOtherTests, fmt.Sprintf("Invalid -> TestRegisterHandlerDuplicateUser"))
		testResult.OtherTestsCases = append(testResult.OtherTestsCases, invalidOtherTests)
	} else {
		testResult.PassedTests++
	}

	t.Log(utils.OrangeColor.InitColor + "\n\nTEST 2 - Duplicate User Registration !!!\n" + utils.EndColor)
	printTestResults(t, codeRes, bodyRes)

	testResult.GeneralTests++
	testResult.TotalTests++
}

func TestLoginHandlerStandard(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Assuming you have a valid login request body with a token
	requestBody := StandardLoginRequestBody

	req, err := http.NewRequest("POST", "localhost:8080/v1/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to capture the handler's response
	rr := httptest.NewRecorder()

	// CALL TESTED FUNCTION !!!
	LoginHandler(rr, req)

	// Verify the responses...
	codeRes := assertResponseStatusCode(t, rr, http.StatusOK)

	if !codeRes {
		invalidOtherTests = append(invalidOtherTests, fmt.Sprintf("Invalid -> TestLoginHandlerStandard"))
		testResult.OtherTestsCases = append(testResult.OtherTestsCases, invalidOtherTests)
	} else {
		testResult.PassedTests++
	}

	t.Log(utils.OrangeColor.InitColor + "\n\nTEST 2 - Standard User Login !!!\n" + utils.EndColor)
	printTestResults(t, codeRes, codeRes) // OBS BodyRes is not realy passed because it is a JWT token, so cannot be expected

	testResult.TotalTests++
	testResult.GeneralTests++
}

// TEST FOR INVALID FIELDS!!!!!!!!!!
func runInvalidTests(t *testing.T, rb []byte, testType, v1route string, invalidValues []string, expectedResponses ...string) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	totalTests := len(invalidValues)
	passedTests := 0
	invalidTests := make([]string, 0)

	for i, value := range invalidValues {
		t.Run(fmt.Sprintf("Invalid%sCase_%d", testType, i+1), func(t *testing.T) {
			// Gerar o RequestBody com o valor inválido a ser testado (primeiro nome, sobrenome, etc.)
			requestBody := GenerateRequestBodyWithInvalidValue(rb, testType, value)

			req, err := makeRequest(v1route, requestBody)
			if err != nil {
				panic("Cant make the request!")
			}

			// Criar um ResponseRecorder para capturar a resposta do handler
			rr := httptest.NewRecorder()

			// CHAMAR A FUNÇÃO TESTADA !!!
			if v1route == "register" {
				RegisterHandler(rr, req)
			} else if v1route == "login" {
				LoginHandler(rr, req)
			}

			// Verificar as respostas...
			codeRes := assertResponseStatusCode(t, rr, http.StatusBadRequest)
			bodyRes := assertResponseBody(t, rr, expectedResponses...)

			if bodyRes && codeRes {
				passedTests++
			} else {
				invalidTests = append(invalidTests, fmt.Sprintf("Invalid %s format: %s \n", testType, value))
				t.Errorf("Test failed: Invalid %s format: %s", testType, value)
			}

			t.Logf(utils.OrangeColor.InitColor+"\n\nTEST INVALID %s %d - User Registration with invalid %s format: %s !!!\n"+utils.EndColor, testType, i+1, testType, value)
			printTestResults(t, codeRes, bodyRes)
		})
	}
	testResult.TotalTests += totalTests
	testResult.PassedTests += passedTests
	testResult.InvalidTestsCases = append(testResult.InvalidTestsCases, invalidTests)
}

// -------------------------------FOR REGISTER ------------------------
// TEST FIRST NAME
func TestRegisterHandlerForInvalidFirstNames(t *testing.T) {
	invalidFirstNames := []string{
		"",
		"select",
		"Joe*",
		"J",
		"JoeDoeJoeDoeJoeDoeJoeDoeJoeDoe",
	}

	expectedResponses := []string{
		`{"error":"Bad Request: {1 : First name must have 2 to 25 characters - and valid characters-words}"}`,
		// Outras respostas esperadas para o bodyRes, se houver...
	}
	runInvalidTests(t, StandardRequestBody, "FirstName", "register", invalidFirstNames, expectedResponses...)
}

// TEST LAST NAME
func TestRegisterHandlerForInvalidLastName(t *testing.T) {
	invalidLastNames := []string{
		"",
		"select",
		"Doe*",
		"Doeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	}

	expectedResponses := []string{
		`{"error":"Bad Request: {1 : Last name must have 1 to 40 characters - and valid characters-words}"}`,
		// Outras respostas esperadas para o bodyRes, se houver...
	}
	runInvalidTests(t, StandardRequestBody, "LastName", "register", invalidLastNames, expectedResponses...)
}

// TEST EMAIL
func TestRegisterHandlerForInvalidEmail(t *testing.T) {
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

	expectedResponses := []string{
		`{"error":"Bad Request: {1 : Invalid email format}"}`,
		// Outras respostas esperadas para o bodyRes, se houver...
	}
	runInvalidTests(t, StandardRequestBody, "Email", "register", invalidEmails, expectedResponses...)
}

// TEST PASSWORD
func TestRegisterHandlerForInvalidPassword(t *testing.T) {
	invalidPasswords := []string{
		"",
		"Ab@1",
		"ThisIsAnExtremelyLongPassword1234567890!@#",
		"password; DROP TABLE users;--",
		"ThisPassContainsSpaces ",
		" MyPassHasLeadingSpaces",
		"MyPassHasInvalidSymbols\000\001",
		"MyPassHasInvalidSymbols\\n\\t\\r",
		"ThisPassIsTooLongItExceedsTheMaxAllowedLimitOfThirtyCharacters",
	}

	expectedResponses := []string{
		`{"error":"Bad Request: {1 : Password must have at least 5 characters - and valid characters-words}"}`,
	}
	runInvalidTests(t, StandardRequestBody, "Password", "register", invalidPasswords, expectedResponses...)
}

// TEST DATE OF BIRTH
func TestRegisterHandlerForInvalidDateOfBirth(t *testing.T) {
	invalidDateOfBirths := []string{
		"",
		"01/01/1990",
		"1990-01",
		"1990-01-01 12:34:56",
		"1990/01/01",
		"1990-13-01",
		"2050-01-01",
		"-1990-01-01",
		"1990-00-01",
		"1990-01-00",
		"1990-01-32",
	}

	expectedResponses := []string{
		`{"error":"Bad Request: {1 : Invalid date of birth}"}`,
		`{"error":"Bad Request: {1 : Date of birth cannot be in the future}"}`,
		`{"error":"Bad Request: Error during signup: invalid date of birth format. It must be in the format YYYY-MM-DD"}`,
		// Outras respostas esperadas para o bodyRes, se houver...
	}
	runInvalidTests(t, StandardRequestBody, "DateOfBirth", "register", invalidDateOfBirths, expectedResponses...)
}

// TEST FOR NICKNAME
func TestRegisterHandlerForInvalidNiNickNames(t *testing.T) {
	invalidNickNames := []string{
		"select",
		"Doe*",
		"J",
		"JoeDoeJoeDoeJoeDoeJoeDoeJoeDoeDoeDoeDoeDoeDoeDoe",
		"J\n",
	}

	expectedResponses := []string{
		`{"error":"Bad Request: {1 : Nickname must have 2 to 30 characters - and valid characters-words}"}`,
		// Outras respostas esperadas para o bodyRes, se houver...
	}
	runInvalidTests(t, StandardRequestBody, "Nickname", "register", invalidNickNames, expectedResponses...)
}

// TEST FIELD
func TestRegisterHandlerForInvalidField(t *testing.T) {
	invalidFields := []string{
		"",
		" ",
		"Software Engineering\n",
		"Software Engineering Software Engineering Software Engineering Software Engineering Software Engineering Software Engineering Software Engineering Software Engineering Software Engineering Software Engineering Software Engineering Software Engineering",
	}

	expectedResponses := []string{
		`{"error":"Bad Request: {1 : Field must have 2 to 50 characters - and valid characters-words}"}`,
	}
	runInvalidTests(t, StandardRequestBody, "Field", "register", invalidFields, expectedResponses...)
}

// TEST BIOGRAPHY
func TestRegisterHandlerForInvalidBiography(t *testing.T) {
	invalidBiographies := []string{
		"",
		"A select.\n",
		"A ",
		"Im a Engineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering SoftwareEngineering Software",
	}

	expectedResponses := []string{
		`{"error":"Bad Request: {1 : Biography must have 3 to 500 characters - and valid characters-words}"}`,
	}
	runInvalidTests(t, StandardRequestBody, "Biography", "register", invalidBiographies, expectedResponses...)
}

// -------------------------------FOR LOGIN ------------------------

func TestLoginHandlerForInvalidEmail(t *testing.T) {
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

	expectedResponses := []string{
		`{"error":"Bad Request: Invalid email format"}`,
		`{"error":"Bad Request: Error during login: User not found -  invalid email or password!"}`,
		// Outras respostas esperadas para o bodyRes, se houver...
	}
	runInvalidTests(t, StandardLoginRequestBody, "Email", "login", invalidEmails, expectedResponses...)
}

// TEST PASSWORD
func TestLoginHandlerForInvalidPassword(t *testing.T) {
	invalidPasswords := []string{
		"",
		"Ab@1",
		"ThisIsAnExtremelyLongPassword1234567890!@#",
		"password; DROP TABLE users;--",
		"ThisPassContainsSpaces ",
		" MyPassHasLeadingSpaces",
		"MyPassHasInvalidSymbols\000\001",
		"MyPassHasInvalidSymbols\\n\\t\\r",
		"ThisPassIsTooLongItExceedsTheMaxAllowedLimitOfThirtyCharacters",
	}

	expectedResponses := []string{
		`{"error":"Bad Request: Password must have at least 5 characters - and valid characters-words"}`,
		`{"error":"Bad Request: Error during login: User not found -  invalid email or password!"}`,
	}
	runInvalidTests(t, StandardLoginRequestBody, "Password", "login", invalidPasswords, expectedResponses...)
}

// PRINT RESULT

func TestPrintTestSummary(t *testing.T) {
	cB := utils.BlueColor.InitColor
	cG := utils.GreenColor.InitColor
	cY := utils.YellowColor.InitColor
	cC := utils.CyanColor.InitColor
	rst := utils.EndColor

	t.Helper()
	t.Logf("\n\nTEST SUMMARY:\n")
	t.Logf("%sTotal Tests: %d%s\n", cB, testResult.TotalTests, rst)
	t.Logf("%sPassed Tests: %d%s\n", cG, testResult.PassedTests, rst)
	t.Logf("%sGeneral Tests: %d%s\n", cC, testResult.GeneralTests, rst)
	t.Logf("%sNumber of fields tested for invalid cases: %d%s\n", cY, len(testResult.InvalidTestsCases), rst)

	printFailedTestsResults(testResult.OtherTestsCases, t)
	printFailedTestsResults(testResult.InvalidTestsCases, t)
}
