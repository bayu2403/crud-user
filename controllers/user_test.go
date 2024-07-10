package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"crud/user/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserTestSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	r    *gin.Engine
}

func (suite *UserTestSuite) SetupTest() {
	var err error
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	suite.DB, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(suite.T(), err)

	suite.mock = mock

	gin.SetMode(gin.TestMode)
	suite.r = gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("email", ValidateEmail)
		v.RegisterValidation("e164", ValidatePhoneNumber)
	}
	models.DB = suite.DB
}

func (suite *UserTestSuite) TearDownTest() {
	sqlDB, err := suite.DB.DB()
	assert.NoError(suite.T(), err)
	sqlDB.Close()
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) TestFindUsers() {
	rows := sqlmock.NewRows([]string{"id", "name", "email", "address", "age", "phone_number", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "John Doe", "john@example.com", "Address 1", 30, "+1234567890", time.Now(), time.Now(), nil)

	suite.mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE deleted_at is null ORDER BY created_at desc$").
		WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/v1/users", nil)
	w := httptest.NewRecorder()
	suite.r.GET("/v1/users", FindUsers)
	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserTestSuite) TestFindUser() {
	rows := sqlmock.NewRows([]string{"id", "name", "email", "address", "age", "phone_number", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "John Doe", "john@example.com", "Address 1", 30, "+1234567890", time.Now(), time.Now(), nil)

	suite.mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE id = \\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
		WithArgs("1", 1).
		WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/v1/users/1", nil)
	w := httptest.NewRecorder()
	suite.r.GET("/v1/users/:id", FindUser)
	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserTestSuite) TestCreateUsers() {
	// Mock input JSON
	input := CreateUserInput{
		Name:        "test",
		Email:       "test@gmail.com",
		Address:     "jalan 123",
		Age:         24,
		PhoneNumber: "+62234567890",
	}
	inputJSON, _ := json.Marshal(input)

	// Mock DB interaction
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(sqlmock.AnyArg(), input.Name, input.Email, input.Address, input.Age, input.PhoneNumber, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()

	// Create request and response recorder
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Route and handle request
	suite.r.POST("/v1/users", CreateUsers)
	suite.r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), input.Name, response["data"].Name)
	assert.Equal(suite.T(), input.Email, response["data"].Email)
	assert.Equal(suite.T(), input.Address, response["data"].Address)
	assert.Equal(suite.T(), input.Age, response["data"].Age)
	assert.Equal(suite.T(), input.PhoneNumber, response["data"].PhoneNumber)
}

func (suite *UserTestSuite) TestUpdateUser() {
	// Mock existing user
	existingUser := models.User{
		ID:          1,
		Name:        "test",
		Email:       "test@gmail.com",
		Address:     "jalan 123",
		Age:         24,
		PhoneNumber: "+62234567890",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE id = \\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
		WithArgs("1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "address", "age", "phone_number", "created_at", "updated_at", "deleted_at"}).
			AddRow(existingUser.ID, existingUser.Name, existingUser.Email, existingUser.Address, existingUser.Age, existingUser.PhoneNumber, existingUser.CreatedAt, existingUser.UpdatedAt, nil))

	// Mock input JSON
	input := UpdateUserInput{
		Name:    "test 2",
		Address: "jalan 1234",
	}
	inputJSON, _ := json.Marshal(input)

	// Mock DB update
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "users" SET (.+) WHERE "id" = \\?`).
		WithArgs(sqlmock.AnyArg(), input.Name, input.Address, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	// Create request and response recorder
	req, _ := http.NewRequest("PATCH", "/v1/users/1", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Route and handle request
	suite.r.PATCH("/v1/users/:id", UpdateUser)
	suite.r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), input.Name, response["data"].Name)
	assert.Equal(suite.T(), input.Address, response["data"].Address)
	assert.Equal(suite.T(), existingUser.Email, response["data"].Email)             // Email should remain unchanged
	assert.Equal(suite.T(), existingUser.PhoneNumber, response["data"].PhoneNumber) // PhoneNumber should remain unchanged
	assert.Equal(suite.T(), existingUser.Age, response["data"].Age)                 // Age should remain unchanged
}

func (suite *UserTestSuite) TestDeleteUser() {
	rows := sqlmock.NewRows([]string{"id", "name", "email", "address", "age", "phone_number", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "John Doe", "john@example.com", "Address 1", 30, "+1234567890", time.Now(), time.Now(), nil)

	suite.mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE id = \\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
		WithArgs("1", 1).
		WillReturnRows(rows)

	req, _ := http.NewRequest("DELETE", "/v1/users/1", nil)
	w := httptest.NewRecorder()
	suite.r.DELETE("/v1/users/:id", DeleteUser)
	suite.r.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserTestSuite) TestFindUserNotFound() {
	// Mock the database query to return an error
	suite.mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE id = \\$1").
		WithArgs("1").
		WillReturnError(gorm.ErrRecordNotFound)

	// Create request and response recorder
	req, _ := http.NewRequest("GET", "/v1/users/1", nil)
	w := httptest.NewRecorder()

	// Route and handle request
	suite.r.GET("/v1/users/:id", FindUser)
	suite.r.ServeHTTP(w, req)

	fmt.Println(w.Body)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *UserTestSuite) TestCreateUsersInvalidInput() {
	// Mock input JSON
	input := CreateUserInput{
		Name:        "test",
		Email:       "wrong-email",
		Address:     "jalan 123",
		Age:         24,
		PhoneNumber: "+62234567890",
	}
	inputJSON, _ := json.Marshal(input)

	// Create request and response recorder
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Route and handle request
	suite.r.POST("/v1/users", CreateUsers)
	suite.r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *UserTestSuite) TestUpdateUserInvalidInput() {
	// Mock DB query for existing user
	existingUser := models.User{
		ID:          1,
		Name:        "test",
		Email:       "test@gmail.com",
		Address:     "jalan 123",
		Age:         24,
		PhoneNumber: "+62234567890",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE id = \\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
		WithArgs("1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "address", "age", "phone_number", "created_at", "updated_at", "deleted_at"}).
			AddRow(existingUser.ID, existingUser.Name, existingUser.Email, existingUser.Address, existingUser.Age, existingUser.PhoneNumber, existingUser.CreatedAt, existingUser.UpdatedAt, nil))

	// Mock invalid input JSON
	input := UpdateUserInput{
		Name:        "test",
		Email:       "wrong-email",
		Address:     "jalan 123",
		Age:         24,
		PhoneNumber: "+62234567890",
	}
	inputJSON, _ := json.Marshal(input)

	// Create request and response recorder
	req, _ := http.NewRequest("PATCH", "/v1/users/1", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Route and handle request
	suite.r.PATCH("/v1/users/:id", UpdateUser)
	suite.r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *UserTestSuite) TestDeleteUserNotFound() {
	// Mock DB query for user not found
	suite.mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE id = \\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
		WithArgs("999", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Create request and response recorder
	req, _ := http.NewRequest("DELETE", "/v1/users/999", nil)
	w := httptest.NewRecorder()

	// Route and handle request
	suite.r.DELETE("/v1/users/:id", DeleteUser)
	suite.r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func TestValidateCreateInput(t *testing.T) {
	type testData struct {
		Input    CreateUserInput
		Expected int
	}

	tests := []testData{
		{Input: CreateUserInput{Name: "A", Address: "Valid Address"}, Expected: http.StatusBadRequest}, // Name too short
		{Input: CreateUserInput{Name: "Valid Name", Address: "A"}, Expected: http.StatusBadRequest},    // Address too short
		{Input: CreateUserInput{Name: "Valid Name", Address: "Valid Address"}, Expected: http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.Input.Name+"_"+test.Input.Address, func(t *testing.T) {
			// Mock Gin context
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			// Invoke validation function
			validateCreateInput(c, &test.Input)

			// Check if error is expected
			assert.Equal(t, test.Expected, c.Writer.Status())
		})
	}
}

func TestValidateUpdateInput(t *testing.T) {
	type testData struct {
		Input    UpdateUserInput
		Expected bool
	}

	tests := []testData{
		{Input: UpdateUserInput{Name: "A", Address: "Valid Address"}, Expected: false}, // Name too short
		{Input: UpdateUserInput{Name: "Valid Name", Address: "A"}, Expected: false},    // Address too short
		{Input: UpdateUserInput{Name: "Valid Name", Address: "Valid Address"}, Expected: true},
		{Input: UpdateUserInput{Name: "", Address: "Valid Address"}, Expected: true}, // No validation error for empty name
	}

	for _, test := range tests {
		t.Run(test.Input.Name+"_"+test.Input.Address, func(t *testing.T) {
			// Mock Gin context
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			validateUpdateInput(c, &test.Input)

			// Check if error is expected
			if test.Expected {
				assert.Equal(t, http.StatusOK, c.Writer.Status())
			} else {
				assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
			}
		})
	}
}
