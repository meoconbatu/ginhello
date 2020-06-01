package main

import (
	handler "ginhello/handlers"
	model "ginhello/models"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var router *gin.Engine

func main() {
	port := os.Getenv("PORT")

	var err error
	db, err := model.NewDB("./models/gorm.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Article{}, &model.User{}, &model.VerificationToken{})

	env := &handler.Env{DB: db}

	router = handler.SetupRouter(env)

	router.RunTLS(":"+port, "server.crt", "server.key")
}
