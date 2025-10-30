package main

import (
	"cursed_backend/db"
	"cursed_backend/router"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env not found")
	}
	log.Println("✅ .env loaded")
	db.InitDB()
	log.Println("✅ DB connected")
	r := router.SetupRouter()
	log.Println("✅ Router setup")
	port := os.Getenv("PORT")
	errorRun := r.Run(":" + port)
	if errorRun != nil {
		log.Println("❌ Server error:", errorRun)
		return
	}
}
