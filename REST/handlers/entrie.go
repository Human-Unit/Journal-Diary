package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"Gin/models"
	pb "Gin/proto/gen/go"

	"github.com/gin-gonic/gin"
    "strconv"
    "io"
)

// This will hold a persistent gRPC client

func GetEntrieData(c *gin.Context) (models.Entry, error) {
    var userLog models.Entry
    if err := c.ShouldBindJSON(&userLog); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid input format",
            "details": err.Error(),
        })
        return userLog, err
    }
    return userLog, nil
}

func CreateEntrie(c *gin.Context) {
    start := time.Now()

    data, err := GetEntrieData(c)
    if err != nil {
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    callStart := time.Now()
    response, err := grpcClient.SaveEntry(ctx, &pb.EntryRequest{
        Title:    data.Title,
        Userid: uint64(data.UserID),

    })
    if err != nil {
        log.Printf("gRPC call failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to communicate with auth service",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "Title": response.Title,
		"UserID": response.Userid,
		"ID": response.Id,
    })

    log.Printf("gRPC call took: %v", time.Since(callStart))
    log.Printf("Total request took: %v", time.Since(start))
}

func DeleteEntrie(c *gin.Context) {
    idStr := c.Param("id")

    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    callStart := time.Now()
    response, err := grpcClient.DeleteEntry(ctx, &pb.DeleteEntryRequest{
        Id: id,
    })
    if err != nil {
        log.Printf("gRPC call failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to communicate with auth service",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": response.Message,
    })

    log.Printf("gRPC call took: %v", time.Since(callStart))
}

func GetAllEntries(c *gin.Context) {
    // Get user ID from query parameter or context
    userIDStr := c.Query("user_id")
    if userIDStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user_id parameter is required"})
        return
    }

    userID, err := strconv.ParseUint(userIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    callStart := time.Now()
    
    stream, err := grpcClient.GetAllEntries(ctx, &pb.GetAllEntriesRequest{
        Userid: userID,
    })
    if err != nil {
        log.Printf("gRPC call failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with auth service"})
        return
    }

    var entries []models.EntryResponse
    
    for {
        entryResponse, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Error receiving from stream: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to receive data from auth service"})
            return
        }

        entries = append(entries, models.EntryResponse{
            ID:      entryResponse.Id,
            Title:   entryResponse.Title,
            UserID:  uint64(entryResponse.Userid),
            Content: entryResponse.Content,
        })
    }

    c.JSON(http.StatusOK, entries)
    log.Printf("gRPC call took: %v", time.Since(callStart))
}

func UpdateEntrie(c *gin.Context) {
    idStr := c.Param("id")

    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    var data models.Entry

    if err := c.ShouldBindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid input format",
            "details": err.Error(),
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    callStart := time.Now()
    response, err := grpcClient.UpdateEntry(ctx, &pb.UpdateEntryRequest{
        Id:       id,
        Title:    data.Title,
        Userid: uint64(data.UserID),
    })
    if err != nil {
        log.Printf("gRPC call failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to communicate with auth service",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": response.Title,
    })

    log.Printf("gRPC call took: %v", time.Since(callStart))
}

type EntryResponse struct {
    ID      uint64 `json:"id"`
    UserID  uint64 `json:"user_id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

func GetEntrie(c *gin.Context) {
    idStr := c.Param("id")
    
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    response, err := grpcClient.GetEntry(ctx, &pb.GetEntryRequest{Id: id})
    if err != nil {
        log.Printf("gRPC call failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with auth service"})
        return
    }

    entry := EntryResponse{
        ID:      response.Id,
        UserID:  response.Userid,
        Title:   response.Title,
        Content: response.Content,
    }

    c.JSON(http.StatusOK, entry)
}