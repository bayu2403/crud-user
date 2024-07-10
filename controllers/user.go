package controllers

import (
	"crud/user/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/swag/example/celler/httputil"
)

type CreateUserInput struct {
	Name string `json:"name" binding:"required" example:"testName"`
	// Check if it's email
	Email   string `json:"email" binding:"required,email" example:"testName@gmail.com"`
	Address string `json:"address" binding:"required" example:"purworejo, jawa tengah, indonesia"`
	Age     int8   `json:"age" binding:"required" example:"24"`
	// Check if it's phoneNumber
	PhoneNumber string `json:"phoneNumber" binding:"required,e164" example:"+6285155678965"`
}

type UpdateUserInput struct {
	Name string `json:"name" example:"testName"`
	// Check if it's email
	Email   string `json:"email" binding:"email" example:"testName@gmail.com"`
	Address string `json:"address" example:"purworejo, jawa tengah, indonesia"`
	Age     int8   `json:"age" example:"24"`
	// Check if it's phoneNumber
	PhoneNumber string `json:"phoneNumber" binding:"e164" example:"+6285155678965"`
}

type CustomError struct {
	Code    int
	Message string
}

func (c *CustomError) Error() string {
	return c.Message
}

// FindUsers godoc
// @Summary      Find All User where not deleted and sorted by created_at
// @Description  find all user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.User
// @Router       /v1/users [get]
func FindUsers(c *gin.Context) {
	var users []models.User
	models.DB.Where("deleted_at is null").Order("created_at desc").Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// ShowAccount godoc
// @Summary      Find by id
// @Description  get by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      404  {object}  httputil.HTTPError
// @Router       /v1/users/{id} [get]
func FindUser(c *gin.Context) {
	var user models.User

	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		httputil.NewError(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// ShowAccount godoc
// @Summary      Create user
// @Description  create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 			 request body controllers.CreateUserInput true "body"
// @Success      200  {object}  models.User
// @Failure      400  {object}  httputil.HTTPError
// @Router       /v1/users [post]
func CreateUsers(c *gin.Context) {
	// Validate input
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	validateCreateInput(c, &input)

	// Create user
	user := models.User{
		Name:        input.Name,
		Email:       input.Email,
		Address:     input.Address,
		Age:         input.Age,
		PhoneNumber: input.PhoneNumber,
		CreatedAt:   time.Now(),
	}
	models.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// ShowAccount godoc
// @Summary      Update user
// @Description  update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param 			 request body controllers.UpdateUserInput true "body"
// @Success      200  {object}  models.User
// @Failure      400  {object}  httputil.HTTPError
// @Router       /v1/users/{id} [patch]
func UpdateUser(c *gin.Context) {
	// Get User if exist
	var user models.User
	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	// Validate input
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	validateUpdateInput(c, &input)
	user.UpdatedAt = time.Now()

	// Update user
	models.DB.Model(&user).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// ShowAccount godoc
// @Summary      Delete user
// @Description  delete user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  httputil.HTTPError
// @Router       /v1/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	// Get model if exist
	var user models.User
	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// Custom validation function for email
// Only check when email not empty
func ValidateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return true // Skip validation if the email is empty
	}
	err := validator.New().Var(email, "email")
	return err == nil
}

// Custom validation function for phone number
// Only check when phone number not empty
func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	if phoneNumber == "" {
		return true // Skip validation if the phone number is empty
	}
	err := validator.New().Var(phoneNumber, "e164")
	return err == nil
}

func validateCreateInput(c *gin.Context, input *CreateUserInput) {
	validations := map[string]string{
		input.Name:    "Name should be more than 1 char",
		input.Address: "Address should be more than 1 char",
	}

	for field, errorMessage := range validations {
		if len(field) < 2 {
			httputil.NewError(c, http.StatusBadRequest, &CustomError{
				Code:    http.StatusBadRequest,
				Message: errorMessage,
			})
			return
		}
	}
}

func validateUpdateInput(c *gin.Context, input *UpdateUserInput) {
	validations := map[string]string{
		input.Name:    "Name should be more than 1 char",
		input.Address: "Address should be more than 1 char",
	}

	for field, errorMessage := range validations {
		if len(field) > 0 && len(field) < 2 {
			httputil.NewError(c, http.StatusBadRequest, &CustomError{
				Code:    http.StatusBadRequest,
				Message: errorMessage,
			})
			return
		}
	}
}
