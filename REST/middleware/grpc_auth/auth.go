package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	pb "Gin/proto/gen/go" // Ensure this matches your actual generated protobuf package
)

// This will hold a persistent gRPC client
var grpcClient pb.AuthServiceClient
type UserLog struct {
	Email string `json:"email"`
	Password string `json:"password"` 

}
func Auth() gin.HandlerFunc{
	return func(c *gin.Context){
		start := time.Now()
		token, err := c.Cookie("token")
		if err != nil {
            log.Printf("Auth failed - no token cookie: %v", err)
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error":   "Unauthorized",
                "message": "Missing authentication token",
            })
            return
        }
		if token == "" {
            log.Println("Auth failed - empty token")
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token format",
            })
            return
        }
		
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // Increased timeout
        defer cancel()

		callStart := time.Now()
		response, err := grpcClient.Auth(ctx, &pb.AuthRequest{
            Token: token,
        })
		if err != nil{
			log.Printf("grpc call failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to commutincate with server",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Status": response.Status,
		})
		log.Printf("gRPC call took: %v", time.Since(callStart))
		log.Printf("Total request took: %v", time.Since(start))
}
}

func GetData(c *gin.Context) (UserLog, error) {
	var UserLog UserLog
	if err := c.ShouldBindJSON(&UserLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input format",
			"details": err.Error(),
		})
		return UserLog, err
	}
	return UserLog, nil
}


func Login(c *gin.Context) {
	// Get token from request
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil { log.Fatalf("failed to connect: %v", err) }
	grpcClient = pb.NewAuthServiceClient(conn)
	data, err := GetData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input format",
			"details": err.Error(),
		})
		return // Error response already handled in GetData
	}
	if grpcClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Authentication service unavailable",
		})
		log.Println("gRPC client not initialized")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	response, err := grpcClient.SendUserLogData(ctx, &pb.UserLogDataRequest{
		Email: data.Email,
		Password: data.Password, // Pass the extracted token string
	})
	
	if err != nil {
		log.Printf("gRPC call failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server failure",
			"details": err.Error(),
		})
		return
	}
	if response == nil || response.Token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid authentication response",
		})
		return
	}
	SetTokenCookie(c, response.Token)
	c.JSON(http.StatusOK, gin.H{
		"Token": response.Token,
	})
}
func SetTokenCookie(c *gin.Context, token string) {
	if token == "" {
		log.Println("Warning: Attempted to set empty token cookie")
		return
	}
	c.SetCookie("token", token, 3600, "/", "", false, true)
}