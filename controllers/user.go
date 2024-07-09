package controllers

import (
	"crud/user/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateUserInput struct {
	Name string `json:"name" binding:"required"`
	// Check if it's email
	Email   string `json:"email" binding:"required,email"`
	Address string `json:"address" binding:"required"`
	Age     int8   `json:"age" binding:"required"`
	// Check if it's phoneNumber
	PhoneNumber string `json:"phoneNumber" binding:"required,e164"`
}

type UpdateUserInput struct {
	Name string `json:"name"`
	// Check if it's email
	Email   string `json:"email" binding:"email"`
	Address string `json:"address"`
	Age     int8   `json:"age"`
	// Check if it's phoneNumber
	PhoneNumber string    `json:"phoneNumber" binding:"e164"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// FindUsers godoc
// @Summary      Find All User where not deleted and sorted by created_at
// @Description  find all user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Router       /users [get]
func FindUsers(c *gin.Context) {
	var users []models.User
	models.DB.Where("deleted_at is null").Order("created_at desc").Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func FindUser(c *gin.Context) {
	var user models.User

	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func CreateUsers(c *gin.Context) {
	// Validate input
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func UpdateUser(c *gin.Context) {
	// Get User if exist
	var user models.User
	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validateUpdateInput(c, &input)
	input.UpdatedAt = time.Now()

	// Update user
	models.DB.Model(&user).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

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

func validateCreateInput(c *gin.Context, input *CreateUserInput) {
	validations := map[string]string{
		input.Name:    "Name should be more than 1 char",
		input.Address: "Address should be more than 1 char",
	}

	for field, errorMessage := range validations {
		if len(field) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
			return
		}
	}
}
