package main

import (
	"Auth/proto/gen/go"
	"context"
	"log"
	"net"
	"google.golang.org/grpc"
	"Auth/auth"
	"Auth/database"
	"gorm.io/gorm"
	"Auth/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//"time"
)
type server struct {
	auth.UnimplementedAuthServiceServer
}



type fastcr struct{
	Email    string `gorm:"uniqueIndex"`
    Password string
}

func (s *server) SendUserLogData(ctx context.Context, req *auth.UserLogDataRequest) (*auth.UserLogDataResponse, error) {
	log.Print("Received string:", req.GetEmail())
	userData := fastcr{
        Email:    req.GetEmail(),
        Password: req.GetPassword(), // Note: In real implementation, password should be hashed
    }
    var dbUser models.LogData
    db := database.GetDB()
    if db == nil {
        log.Println("Database connection is nil")
        return nil, status.Error(codes.Internal, "database connection error")
    }
    result := db.WithContext(ctx).Where("email = ?", userData.Email).First(&dbUser)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        log.Printf("Failed to retrieve user: %v", result.Error)
        return nil, status.Error(codes.Internal, "failed to retrieve user")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userData.Password)); err != nil {
        log.Println("Password verification failed:", err)
        return nil, status.Error(codes.Unauthenticated, "invalid credentials")
    }
	token, err := middleware.CreateToken(middleware.UserLog(userData))
	if err != nil{
		status.Error(404, "blya")
	}
	return &auth.UserLogDataResponse{
		Token :token,
	}, nil
}

func (s *server) SaveUser(ctx context.Context, req *auth.SaveUserRequest) (*auth.SaveUserResponse, error) {
    // 1. Input validation
	
	if req.Email == "" || req.Password == "" {
        return nil, status.Error(codes.InvalidArgument, "email and password are required")
    }

    // 2. Password hashing
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
    if err != nil {
        log.Printf("Failed to hash password: %v", err)
        return nil, status.Error(codes.Internal, "failed to process password")
    }

    // 3. Prepare user data
    userData := models.LogData{
        Email:    req.Email,
        Password: string(hash), // Store hashed password, not plaintext
    }

    // 4. Database operation 	
    db := database.GetDB()
    if db == nil {
        log.Println("Database connection is nil")
        return nil, status.Error(codes.Internal, "database connection error")
    }
    result := db.WithContext(ctx).Create(&userData)
    if result.Error != nil {
        log.Printf("Failed to create user: %v", result.Error)
        return nil, status.Error(codes.Internal, "failed to create user")
    }
    // 5. Success response
    return &auth.SaveUserResponse{
        Message: "User created successfully",// Assuming your model has an ID field
    }, result.Error
}
func (s *server) User(ctx context.Context, req *auth.UserRequest) (*auth.UserResponse, error) {
    if req.Email == "" {
        return nil, status.Error(codes.InvalidArgument, "email is required")
    }

    db := database.GetDB()
    if db == nil {
        log.Println("Database connection is nil")
        return nil, status.Error(codes.Internal, "database connection error")
    }
    var user models.LogData
    result := db.WithContext(ctx).Where("email = ?", req.Email).First(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        log.Printf("Failed to retrieve user: %v", result.Error)
        return nil, status.Error(codes.Internal, "failed to retrieve user")
    }

    // 3. Success response
    return &auth.UserResponse{
        Response: user.Email + user.CreatedAt.String(),
    }, nil
}

func (s *server) SaveEntry(ctx context.Context, req *auth.EntryRequest) (*auth.EntryResponse, error){
    entry := models.Entry{
        Title: req.Title,
        Content : req.Content,
        UserID : uint(req.Userid),
    }
    db := database.GetDB()
    if condition := db == nil; condition {
        log.Println("Database connection is nil")
        return nil, status.Error(codes.Internal, "database connection error")
    }
    result := db.WithContext(ctx).Create(&entry)
    if result.Error != nil {
        log.Printf("Failed to create entry: %v", result.Error)
        return nil, status.Error(codes.Internal, "failed to create entry")
    }
    return &auth.EntryResponse{
        Userid: uint64(entry.UserID),
        Title: entry.Title,
        Id: uint64(entry.ID),
    }, result.Error
}

func (s *server) UpdateEntry(ctx context.Context, req *auth.UpdateEntryRequest) (*auth.UpdateEntryResponse, error) {
    if req.Id == 0 || req.Title == "" || req.Content == "" {
        return nil, status.Error(codes.InvalidArgument, "entry ID, title, and content are required")
    }
    db := database.GetDB()
    if db == nil {
        log.Println("Database connection is nil")
        return nil, status.Error(codes.Internal, "database connection error")
    }
    var entry models.Entry
    result := db.WithContext(ctx).Where("id = ?", req.Id).First(&entry)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, status.Error(codes.NotFound, "entry not found")
        }
        log.Printf("Failed to retrieve entry: %v", result.Error)
        return nil, status.Error(codes.Internal, "failed to retrieve entry")
    }
    // Update the entry
    entry.Title = req.Title
    entry.Content = req.Content
    if err := db.WithContext(ctx).Save(&entry).Error; err != nil {
        log.Printf("Failed to update entry: %v", err)
        return nil, status.Error(codes.Internal, "failed to update entry")
    }
    return &auth.UpdateEntryResponse{
        Title: "Entry updated successfully",
        Id: uint64(entry.ID),
    }, nil
}

func (s *server) Auth(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error){
    middleware.Auth()
    if req.Token == "" {
        return nil, status.Error(codes.InvalidArgument, "token is required")
    }
    return &auth.AuthResponse{
        Status: "Auth service is running" + req.Token,
    }, nil
}
func (s *server) DeleteEntry(ctx context.Context, req *auth.DeleteEntryRequest) (*auth.DeleteEntryResponse, error){
    if req.Id == 0 {
        return nil, status.Error(codes.InvalidArgument, "entry ID is required")
    }
    db := database.GetDB()
    if db == nil {
        log.Println("Database connection is nil")
        return nil, status.Error(codes.Internal, "database connection error")
    }
    var entry models.Entry
    result := db.WithContext(ctx).Where("id = ?", req.Id).First(&entry)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, status.Error(codes.NotFound, "entry not found")
        }
        log.Printf("Failed to retrieve entry: %v", result.Error)
        return nil, status.Error(codes.Internal, "failed to retrieve entry")
    }
    // Delete the entry
    if err := db.WithContext(ctx).Delete(&entry).Error; err != nil {
        log.Printf("Failed to delete entry: %v", err)
        return nil, status.Error(codes.Internal, "failed to delete entry")
    }
    return &auth.DeleteEntryResponse{
        Message: "Entry deleted successfully",
    }, nil
}

func (s *server) GetEntry(ctx context.Context, req *auth.GetEntryRequest) (*auth.GetEntryResponse, error){
    if req.Id == 0 {
        return nil, status.Error(codes.InvalidArgument, "entry ID is required")
    }
    db := database.GetDB()
    if db == nil {
        log.Println("Database connection is nil")
        return nil, status.Error(codes.Internal, "database connection error")
    }
    var entry models.Entry
    result := db.WithContext(ctx).Where("id = ?", req.Id).First(&entry)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, status.Error(codes.NotFound, "entry not found")
        }
        log.Printf("Failed to retrieve entry: %v", result.Error)
        return nil, status.Error(codes.Internal, "failed to retrieve entry")
    }
    return &auth.GetEntryResponse{
        Id: uint64(entry.ID),
        Title: entry.Title,
        Content: entry.Content,
        Userid: uint64(entry.UserID),
    }, nil
}

func (s *server) GetAllEntries(req *auth.GetAllEntriesRequest, stream auth.AuthService_GetAllEntriesServer) error {
    if req.Userid == 0 {
        return status.Errorf(codes.InvalidArgument, "User ID is required")
    }

    db := database.GetDB()
    if db == nil {
        return status.Errorf(codes.Internal, "Database connection unavailable")
    }

    var entries []models.Entry
    result := db.WithContext(stream.Context()).Where("user_id = ?", req.Userid).Find(&entries)
    if result.Error != nil {
        log.Printf("Failed to retrieve entries: %v", result.Error)
        return status.Errorf(codes.Internal, "Failed to retrieve entries: %v", result.Error)
    }

    // Stream each entry individually
    for _, entry := range entries {
        if err := stream.Send(&auth.EntryResponse{
            Id:      uint64(entry.ID),
            Title:   entry.Title,
            Content: entry.Content,
            Userid:  uint64(entry.UserID),
            // Add other fields as needed
        }); err != nil {
            log.Printf("Failed to send entry: %v", err)
            return status.Errorf(codes.Internal, "Failed to stream entries: %v", err)
        }
    }

    return nil
}

var (
	db *gorm.DB
)

func main() {
    // Connect to the database first
    if err := database.Connect(); err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }

    // Now get DB instance (optional if you always use database.GetDB() directly)
    db = database.GetDB()

    // Create TCP listener
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // Create gRPC server instance
    s := grpc.NewServer()

    // Register your service
    auth.RegisterAuthServiceServer(s, &server{})

    // Start serving
    log.Println("gRPC server listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
