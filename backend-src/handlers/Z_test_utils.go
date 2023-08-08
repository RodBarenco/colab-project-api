// I cant call handlers_test case I need to import a variable that is just in handlers package and cant be exported
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
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

type LoginRequestBody struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
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

var StandardLoginRequestBody = []byte(`{
	"Email": "exemple@example.com",
	"Password": "123456"
}`)

// sets the request body. -------------------------------------------------------------

func GenerateRequestBodyWithInvalidValue(testBody []byte, testType string, invalidValue string) []byte {
	var body TestRequestBody
	err := json.Unmarshal(testBody, &body)

	if err != nil {
		var body2 LoginRequestBody
		err2 := json.Unmarshal(testBody, &body2)
		if err2 != nil {
			panic("Failed to unmarshal request body")
		}
		return generateModifiedLoginRequestBody(body2, testType, invalidValue)
	}

	return generateModifiedTestRequestBody(body, testType, invalidValue)
}

func generateModifiedTestRequestBody(body TestRequestBody, testType string, invalidValue string) []byte {
	switch testType {
	case "FirstName":
		body.FirstName = invalidValue
	case "LastName":
		body.LastName = invalidValue
	case "Email":
		body.Email = invalidValue
	case "Password":
		body.Password = invalidValue
	case "DateOfBirth":
		body.DateOfBirth = invalidValue
	case "Nickname":
		body.Nickname = invalidValue
	case "Field":
		body.Field = invalidValue
	case "Biography":
		body.Biography = invalidValue
	default:
		panic("Invalid testType")
	}

	rb, err := json.Marshal(body)
	if err != nil {
		panic("Failed to marshal modified request body")
	}
	return rb
}

func generateModifiedLoginRequestBody(body LoginRequestBody, testType string, invalidValue string) []byte {
	switch testType {
	case "Email":
		body.Email = invalidValue
	case "Password":
		body.Password = invalidValue
	default:
		panic("Invalid testType")
	}

	rb, err := json.Marshal(body)
	if err != nil {
		panic("Failed to marshal modified request body")
	}
	return rb
}

// MAKE REQUEST V1

func makeRequest(v1req string, requestBody []byte) (*http.Request, error) {
	route := "localhost:8080/v1/"
	switch v1req {
	case "register":
		route += "register"
	case "login":
		route += "login"
	default:
		return nil, fmt.Errorf("Invalid request")
	}

	req, err := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	return req, nil
}

//-------------------------------------------------------------------------------------
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

func assertResponseBody(t *testing.T, rr *httptest.ResponseRecorder, expectedResponses ...string) bool {
	t.Helper()
	c := utils.PurpleColor.InitColor
	e := utils.EndColor

	for _, expectedResponse := range expectedResponses {
		if rr.Body.String() == expectedResponse {
			return true
		}
	}

	t.Errorf("\n%s>>>Expected response body:%s\n %v \n%s>>> BUT GOT:%s\n %v \n", c, e, expectedResponses, c, e, rr.Body.String())
	return false
}

// RESULTS
type TestResult struct {
	TotalTests        int
	PassedTests       int
	GeneralTests      int
	InvalidTestsCases [][]string
	OtherTestsCases   [][]string
}

var testResult TestResult

func allSlicesEmpty(slice [][]string) bool {
	for _, innerSlice := range slice {
		for _, str := range innerSlice {
			if str != "" {
				return false
			}
		}
	}
	return true
}

func printFailedTestsResults(slice [][]string, t *testing.T) {
	if len(slice) > 0 && !allSlicesEmpty(slice) {
		t.Logf("%s\n\nInvalid Test Cases:\n%s", utils.RedColor.InitColor, utils.EndColor)
		for _, testCase := range slice {
			if len(testCase) > 0 {
				t.Logf("- %s%s\n%s", utils.RedColor.InitColor, testCase, utils.EndColor)
			}
		}
	}
}
