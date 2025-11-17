package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	proto_session "github.com/anhvanhoa/sf-proto/gen/session/v1"
	proto_user "github.com/anhvanhoa/sf-proto/gen/user/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var serverAddress string

func init() {
	viper.SetConfigFile("dev.config.yaml")
	viper.ReadInConfig()
	serverAddress = fmt.Sprintf("%s:%s", viper.GetString("host_grpc"), viper.GetString("port_grpc"))
}

type GRPCClient struct {
	userClient    proto_user.UserServiceClient
	sessionClient proto_session.SessionServiceClient
	conn          *grpc.ClientConn
}

func NewGRPCClient(address string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	return &GRPCClient{
		userClient:    proto_user.NewUserServiceClient(conn),
		sessionClient: proto_session.NewSessionServiceClient(conn),
		conn:          conn,
	}, nil
}

func (c *GRPCClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// User Service Tests
func (c *GRPCClient) TestGetUserById() {
	fmt.Println("\n=== Test GetUserById ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userIDStr = strings.TrimSpace(userIDStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.userClient.GetUserById(ctx, &proto_user.GetUserByIdRequest{
		Id: userIDStr,
	})
	if err != nil {
		fmt.Printf("Error calling GetUserById: %v\n", err)
		return
	}

	fmt.Printf("User found:\n")
	fmt.Printf("ID: %s\n", resp.User.Id)
	fmt.Printf("Email: %s\n", resp.User.Email)
	fmt.Printf("Phone: %s\n", resp.User.Phone)
	fmt.Printf("Full Name: %s\n", resp.User.FullName)
	fmt.Printf("Avatar: %s\n", resp.User.Avatar)
	fmt.Printf("Bio: %s\n", resp.User.Bio)
	fmt.Printf("Address: %s\n", resp.User.Address)
	fmt.Printf("Status: %s\n", resp.User.Status)
	fmt.Printf("Created At: %s\n", resp.User.CreatedAt.AsTime().Format(time.RFC3339))
	fmt.Printf("Updated At: %s\n", resp.User.UpdatedAt.AsTime().Format(time.RFC3339))
	fmt.Printf("Birthday: %s\n", resp.User.Birthday.AsTime().Format(time.RFC3339))
}

func (c *GRPCClient) TestUpdateUserById() {
	fmt.Println("\n=== Test UpdateUserById ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter user ID to update: ")
	userIDStr, _ := reader.ReadString('\n')
	userIDStr = strings.TrimSpace(userIDStr)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter phone: ")
	phone, _ := reader.ReadString('\n')
	phone = strings.TrimSpace(phone)

	fmt.Print("Enter full name: ")
	fullName, _ := reader.ReadString('\n')
	fullName = strings.TrimSpace(fullName)

	fmt.Print("Enter avatar URL: ")
	avatar, _ := reader.ReadString('\n')
	avatar = strings.TrimSpace(avatar)

	fmt.Print("Enter bio: ")
	bio, _ := reader.ReadString('\n')
	bio = strings.TrimSpace(bio)

	fmt.Print("Enter address: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	fmt.Print("Enter status (active/inactive): ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSpace(status)

	fmt.Print("Enter birthday (YYYY-MM-DD): ")
	birthdayStr, _ := reader.ReadString('\n')
	birthdayStr = strings.TrimSpace(birthdayStr)

	fmt.Print("Enter role IDs (comma-separated): ")
	roleIDsStr, _ := reader.ReadString('\n')
	roleIDsStr = strings.TrimSpace(roleIDsStr)
	var roleIDs []string
	if roleIDsStr != "" {
		roleIDParts := strings.Split(roleIDsStr, ",")
		for _, part := range roleIDParts {
			part = strings.TrimSpace(part)
			roleIDs = append(roleIDs, part)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var birthday *timestamppb.Timestamp
	if birthdayStr != "" {
		birthdayTime, err := time.Parse(time.RFC3339, birthdayStr)
		if err != nil {
			fmt.Printf("Invalid birthday format: %v\n", err)
			return
		}
		birthday = timestamppb.New(birthdayTime)
	}

	resp, err := c.userClient.UpdateUser(ctx, &proto_user.UpdateUserRequest{
		Id:       userIDStr,
		Email:    email,
		Phone:    phone,
		FullName: fullName,
		Avatar:   avatar,
		Bio:      bio,
		Address:  address,
		Status:   status,
		Birthday: birthday,
		RoleIds:  roleIDs,
	})
	if err != nil {
		fmt.Printf("Error calling UpdateUserById: %v\n", err)
		return
	}

	fmt.Printf("User updated successfully:\n")
	fmt.Printf("ID: %s\n", resp.UserInfo.Id)
	fmt.Printf("Email: %s\n", resp.UserInfo.Email)
	fmt.Printf("Phone: %s\n", resp.UserInfo.Phone)
	fmt.Printf("Full Name: %s\n", resp.UserInfo.FullName)
	fmt.Printf("Avatar: %s\n", resp.UserInfo.Avatar)
	fmt.Printf("Bio: %s\n", resp.UserInfo.Bio)
	fmt.Printf("Address: %s\n", resp.UserInfo.Address)
	fmt.Printf("Birthday: %s\n", resp.UserInfo.Birthday.AsTime().Format(time.RFC3339))
}

// Session Service Tests
func (c *GRPCClient) TestGetSessions() {
	fmt.Println("\n=== Test GetSessions ===")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.sessionClient.GetAllSessions(ctx, &proto_session.GetAllSessionsRequest{})
	if err != nil {
		fmt.Printf("Error calling GetSessions: %v\n", err)
		return
	}

	fmt.Printf("Found %d sessions:\n", len(resp.Sessions))
	for i, session := range resp.Sessions {
		fmt.Printf("Session %d:\n", i+1)
		fmt.Printf("  Token: %s\n", session.Token)
		fmt.Printf("  User ID: %s\n", session.UserId)
		fmt.Printf("  Type: %s\n", session.Type)
		fmt.Printf("  OS: %s\n", session.Os)
		fmt.Printf("  Expired At: %s\n", session.ExpiredAt.AsTime().Format(time.RFC3339))
		fmt.Printf("  Created At: %s\n", session.CreatedAt.AsTime().Format(time.RFC3339))
		fmt.Println()
	}
}

func (c *GRPCClient) TestGetSessionsByUserId() {
	fmt.Println("\n=== Test GetSessionsByUserId ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userIDStr = strings.TrimSpace(userIDStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.sessionClient.GetSessionsByUserId(ctx, &proto_session.GetSessionsByUserIdRequest{
		UserId: userIDStr,
	})
	if err != nil {
		fmt.Printf("Error calling GetSessionsByUserId: %v\n", err)
		return
	}

	fmt.Printf("Found %d sessions for user %s:\n", len(resp.Sessions), userIDStr)
	for i, session := range resp.Sessions {
		fmt.Printf("Session %d:\n", i+1)
		fmt.Printf("  Token: %s\n", session.Token)
		fmt.Printf("  User ID: %s\n", session.UserId)
		fmt.Printf("  Type: %s\n", session.Type)
		fmt.Printf("  OS: %s\n", session.Os)
		fmt.Printf("  Expired At: %s\n", session.ExpiredAt.AsTime().Format(time.RFC3339))
		fmt.Printf("  Created At: %s\n", session.CreatedAt.AsTime().Format(time.RFC3339))
		fmt.Println()
	}
}

func (c *GRPCClient) TestDeleteSessionByTypeAndToken() {
	fmt.Println("\n=== Test DeleteSessionByTypeAndToken ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter session type: ")
	sessionType, _ := reader.ReadString('\n')
	sessionType = strings.TrimSpace(sessionType)

	fmt.Print("Enter token: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.sessionClient.DeleteSessionByTypeAndToken(ctx, &proto_session.DeleteSessionByTypeAndTokenRequest{
		Type:  sessionType,
		Token: token,
	})
	if err != nil {
		fmt.Printf("Error calling DeleteSessionByTypeAndToken: %v\n", err)
		return
	}

	fmt.Printf("Delete result: %s\n", resp.Message)
	fmt.Printf("Success: %t\n", resp.Success)
}

func (c *GRPCClient) TestDeleteSessionByTypeAndUser() {
	fmt.Println("\n=== Test DeleteSessionByTypeAndUser ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter session type: ")
	sessionType, _ := reader.ReadString('\n')
	sessionType = strings.TrimSpace(sessionType)

	fmt.Print("Enter user ID: ")
	userIDStr, _ := reader.ReadString('\n')
	userIDStr = strings.TrimSpace(userIDStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.sessionClient.DeleteSessionByTypeAndUser(ctx, &proto_session.DeleteSessionByTypeAndUserRequest{
		Type:   sessionType,
		UserId: userIDStr,
	})
	if err != nil {
		fmt.Printf("Error calling DeleteSessionByTypeAndUser: %v\n", err)
		return
	}

	fmt.Printf("Delete result: %s\n", resp.Message)
	fmt.Printf("Success: %t\n", resp.Success)
}

func (c *GRPCClient) TestDeleteSessionExpired() {
	fmt.Println("\n=== Test DeleteSessionExpired ===")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.sessionClient.DeleteSessionExpired(ctx, &proto_session.DeleteSessionExpiredRequest{})
	if err != nil {
		fmt.Printf("Error calling DeleteSessionExpired: %v\n", err)
		return
	}

	fmt.Printf("Delete result: %s\n", resp.Message)
	fmt.Printf("Success: %t\n", resp.Success)
}

func printMenu() {
	fmt.Println("\n=== gRPC User Service Test Client ===")
	fmt.Println("1. User Service Tests")
	fmt.Println("  1.1 Get User By ID")
	fmt.Println("  1.2 Lock User By ID (Not implemented)")
	fmt.Println("  1.3 Update User By ID")
	fmt.Println("2. Session Service Tests")
	fmt.Println("  2.1 Get All Sessions")
	fmt.Println("  2.2 Get Sessions By User ID")
	fmt.Println("  2.3 Delete Session By Type And Token")
	fmt.Println("  2.4 Delete Session By Type And User")
	fmt.Println("  2.5 Delete Session Expired")
	fmt.Println("0. Exit")
	fmt.Print("Enter your choice: ")
}

func main() {
	// Get server address from command line or use default
	address := serverAddress
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	fmt.Printf("Connecting to gRPC server at %s...\n", address)
	client, err := NewGRPCClient(address)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected successfully!")

	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1.1":
			client.TestGetUserById()
		case "1.3":
			client.TestUpdateUserById()
		case "2.1":
			client.TestGetSessions()
		case "2.2":
			client.TestGetSessionsByUserId()
		case "2.3":
			client.TestDeleteSessionByTypeAndToken()
		case "2.4":
			client.TestDeleteSessionByTypeAndUser()
		case "2.5":
			client.TestDeleteSessionExpired()
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
