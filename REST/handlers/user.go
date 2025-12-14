package handlers

import (
    "context"
    "log"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    pb "Gin/proto/gen/go"
)

func GetUser(c *gin.Context) {
	email := c.Param("email")
	start := time.Now()
	if email == "" {
        c.JSON(400, gin.H{"error": "Email is required"})
        return
    }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	callStart := time.Now()
	response, err := grpcClient.User(ctx, &pb.UserRequest{
		Email: email,
	})
	if err != nil{
		log.Printf("grpc call failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commutincate with server",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
        "Status": response.Response,
    })
	log.Printf("gRPC call took: %v", time.Since(callStart))
    log.Printf("Total request took: %v", time.Since(start))
}

func DeleteUser(c *gin.Context) {
	email := c.Param("email")
	start := time.Now()
	if email == "" {
        c.JSON(400, gin.H{"error": "Email is required"})
        return
    }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	callStart := time.Now()
	response, err := grpcClient.User(ctx, &pb.UserRequest{
		Email: email,
	})
	if err != nil{
		log.Printf("grpc call failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commitincate with server",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
        "Status": response.Response,
    })
	log.Printf("gRPC call took: %v", time.Since(callStart))
    log.Printf("Total request took: %v", time.Since(start))
}