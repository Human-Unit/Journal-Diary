package main

import (
	"Gin/handlers"
	"log"
	auth "Gin/middleware/grpc_auth"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize gRPC connection once
	if err := handlers.InitGRPC(); err != nil {
        log.Fatalf("Failed to connect to auth service: %v", err)
    }

	router := gin.Default()

	// Recovery middleware
	router.Use(gin.Recovery())

	// Public routes
	router.POST("/register", handlers.CreateUser)
	router.POST("/login", auth.Login)

	user := router.Group("/user")

	user.GET("/GetUser/:email", handlers.GetUser)
	user.GET("/DeleteUser/:email", handlers.DeleteUser)
	
	entrie := router.Group("/entrie")

	entrie.Use(auth.Auth()) 
	{
	entrie.POST("/create", handlers.CreateEntrie)
	entrie.GET("/get/:id", handlers.GetEntrie)
	entrie.GET("/getAll", handlers.GetAllEntries)
	entrie.PUT("/update/:id", handlers.UpdateEntrie)
	entrie.DELETE("/delete/:id", handlers.DeleteEntrie)
	}
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
