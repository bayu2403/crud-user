package models

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}
	// read .env file
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT")) // don't forget to convert int since port is int type.
	user := os.Getenv("USER_DB")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	// Initialize gorm
	var builder strings.Builder
	builder.WriteString("host=")
	builder.WriteString(host)
	builder.WriteString(" user=")
	builder.WriteString(user)
	builder.WriteString(" password=")
	builder.WriteString(pass)
	builder.WriteString(" dbname=")
	builder.WriteString(dbname)
	builder.WriteString(fmt.Sprintf(" port=%d", port))
	builder.WriteString(" sslmode=disable TimeZone=Asia/Jakarta")

	dsn := builder.String()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return
	}

	DB = db
}
