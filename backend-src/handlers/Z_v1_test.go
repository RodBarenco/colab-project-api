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

// TEST FOR INVALID FIELDS!!!!!!!!!!
func runInvalidTests(t *testing.T, testType string, invalidValues []string, expectedResponses ...string) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	totalTests := len(invalidValues)
	passedTests := 0
	invalidTests := make([]string, 0)

	for i, value := range invalidValues {
		t.Run(fmt.Sprintf("Invalid%sCase_%d", testType, i+1), func(t *testing.T) {
			// Gerar o RequestBody com o valor inválido a ser testado (primeiro nome, sobrenome, etc.)
			requestBody := GenerateRequestBodyWithInvalidValue(testType, value)

			req, err := http.NewRequest("POST", "localhost:8080/v1/register", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Criar um ResponseRecorder para capturar a resposta do handler
			rr := httptest.NewRecorder()

			// CHAMAR A FUNÇÃO TESTADA !!!
			RegisterHandler(rr, req)

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
	runInvalidTests(t, "FirstName", invalidFirstNames, expectedResponses...)
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
	runInvalidTests(t, "LastName", invalidLastNames, expectedResponses...)
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
	runInvalidTests(t, "Email", invalidEmails, expectedResponses...)
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
	runInvalidTests(t, "Password", invalidPasswords, expectedResponses...)
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
	runInvalidTests(t, "DateOfBirth", invalidDateOfBirths, expectedResponses...)
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
	runInvalidTests(t, "Nickname", invalidNickNames, expectedResponses...)
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
	runInvalidTests(t, "Field", invalidFields, expectedResponses...)
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
	runInvalidTests(t, "Biography", invalidBiographies, expectedResponses...)
}

// PRINT RESULT

func TestPrintTestSummary(t *testing.T) {
	totalTestsColor := utils.BlueColor.InitColor
	passedTestsColor := utils.GreenColor.InitColor
	invalidTestsColor := utils.YellowColor.InitColor
	generalTestColor := utils.CyanColor.InitColor
	failColor := utils.RedColor.InitColor
	resetColor := utils.EndColor

	t.Helper()
	t.Logf("\n\nTEST SUMMARY:\n")
	t.Logf("%sTotal Tests: %d%s\n", totalTestsColor, testResult.TotalTests, resetColor)
	t.Logf("%sPassed Tests: %d%s\n", passedTestsColor, testResult.PassedTests, resetColor)
	t.Logf("%sGeneral Tests: %d%s\n", generalTestColor, testResult.GeneralTests, resetColor)
	t.Logf("%sNumber of fields tested for invalid cases: %d%s\n", invalidTestsColor, len(testResult.InvalidTestsCases), resetColor)

	if len(testResult.OtherTestsCases) > 0 && !allSlicesEmpty(testResult.OtherTestsCases) {
		t.Logf("%s\n\nFaield General Tests:\n%s", failColor, resetColor)
		for _, testCase := range testResult.OtherTestsCases {
			if len(testCase) > 0 {
				t.Logf("- %s%s\n%s", failColor, testCase, resetColor)
			}
		}
	}

	if len(testResult.InvalidTestsCases) > 0 && !allSlicesEmpty(testResult.InvalidTestsCases) {
		t.Logf("%s\n\nInvalid Test Cases:\n%s", failColor, resetColor)
		for _, testCase := range testResult.InvalidTestsCases {
			if len(testCase) > 0 {
				t.Logf("- %s%s\n%s", failColor, testCase, resetColor)
			}
		}
	}
}
