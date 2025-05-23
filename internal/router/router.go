package router

import (
	"log"
	
	"github.com/roma2099/uril-go/internal/middleware"
	"github.com/roma2099/uril-go/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/contrib/websocket"
)

func SetUpRoutes(app *fiber.App){
	api:= app.Group("/api", logger.New())
	api.Get("/",handler.Hello)

	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// User
	user := api.Group("/user")
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Patch("/", middleware.Protected(), handler.UpdateUserCountry)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// Product
	product := api.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:id", handler.GetProduct)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
	product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)

	// Game
	game :=api.Group("/game")
	game.Get("/history/:user_id", handler.GetGamesHistory)
	game.Get("/:id",handler.GetGame)
	game.Use(func(c *fiber.Ctx) error {
				if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})
	game.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		if err := handler.GameWebSocketHandler(c); err != nil {
			log.Printf("WebSocket handler error: %v", err)
		}
	})) 
}

