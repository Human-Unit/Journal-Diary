package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"Gin/models"
	pb "Gin/proto/gen/go"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// This will hold a persistent gRPC client
var grpcClient pb.AuthServiceClient

func InitGRPC() error {
    conn, err := grpc.Dial(
        "localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(), // Wait for connection to be established
    )
    if err != nil {
        return err
    }

    grpcClient = pb.NewAuthServiceClient(conn)
    log.Println("Successfully connected to gRPC server")
    return nil
}

func GetData(c *gin.Context) (models.UserLog, error) {
    var userLog models.UserLog
    if err := c.ShouldBindJSON(&userLog); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid input format",
            "details": err.Error(),
        })
        return userLog, err
    }
    return userLog, nil
}

func CreateUser(c *gin.Context) {
    start := time.Now()

    data, err := GetData(c)
    if err != nil {
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    callStart := time.Now()
    response, err := grpcClient.SaveUser(ctx, &pb.SaveUserRequest{
        Email:    data.Email,
        Password: data.Password,
    })
    if err != nil {
        log.Printf("gRPC call failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to communicate with auth service",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "Status": response.Message,
    })

    log.Printf("gRPC call took: %v", time.Since(callStart))
    log.Printf("Total request took: %v", time.Since(start))
}
