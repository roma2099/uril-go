package main

import (
	"log"
	"github.com/roma2099/uril-go/internal/database"
	"github.com/roma2099/uril-go/internal/router"

	"github.com/gofiber/fiber/v2"
)

func main(){
	app:= fiber.New(fiber.Config{
		AppName:"Romario testing Apps",
		CaseSensitive: true,
		StrictRouting: false,
		ServerHeader:"Fiber",

	})
	database.ConnectDB()
	router.SetUpRoutes(app)
	log.Fatal(app.Listen(":3000"))
}