// File: main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models"
	"github.com/vahidlotfi71/Task_Manager/Routes"
)

func main() {
	// ۱) بارگذاری env
	if err := Config.Getenv(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime Error: Could not load environment variables: %s\n", err.Error())
		os.Exit(2)
	}
	fmt.Println("env loaded successfully")

	// ۲) اتصال به DB
	if err := Config.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Connection Error: Could not connect to the database\n%v\n", err.Error())
		os.Exit(2)
	}
	fmt.Println("Successfully connected to PostgreSQL database!")

	// ۳) AutoMigrate
	models := []interface{}{
		&Models.Task{},
	}
	for _, m := range models {
		if err := Config.DB.AutoMigrate(m); err != nil {
			log.Fatalf("migrate: %v", err)
		}
	}

	// ۴) راه‌اندازی Gin و روت‌ها
	r := gin.Default()
	Routes.SetupRoutes(r)

	// ۵) اجرای سرور
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
