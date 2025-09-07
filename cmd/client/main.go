package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	proto_role "github.com/anhvanhoa/sf-proto/gen/role/v1"
	proto_session "github.com/anhvanhoa/sf-proto/gen/session/v1"
	proto_user "github.com/anhvanhoa/sf-proto/gen/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	serverAddress = "localhost:50050" // Default gRPC server address
)

type GRPCClient struct {
	userClient    proto_user.UserServiceClient
	sessionClient proto_session.SessionServiceClient
	roleClient    proto_role.RoleServiceClient
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
		roleClient:    proto_role.NewRoleServiceClient(conn),
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

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid user ID: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.userClient.GetUserById(ctx, &proto_user.GetUserByIdRequest{
		Id: strconv.FormatInt(userID, 10),
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

func (c *GRPCClient) TestDeleteUserById() {
	fmt.Println("\n=== Test DeleteUserById ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter user ID to delete: ")
	userIDStr, _ := reader.ReadString('\n')
	userIDStr = strings.TrimSpace(userIDStr)

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid user ID: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.userClient.DeleteUser(ctx, &proto_user.DeleteUserRequest{
		Id: strconv.FormatInt(userID, 10),
	})
	if err != nil {
		fmt.Printf("Error calling DeleteUserById: %v\n", err)
		return
	}

	fmt.Printf("Delete result: %s\n", resp.Message)
	fmt.Printf("Success: %t\n", resp.Success)
}

func (c *GRPCClient) TestUpdateUserById() {
	fmt.Println("\n=== Test UpdateUserById ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter user ID to update: ")
	userIDStr, _ := reader.ReadString('\n')
	userIDStr = strings.TrimSpace(userIDStr)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid user ID: %v\n", err)
		return
	}

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
	birthday, err := time.Parse("2006-01-02", birthdayStr)
	if err != nil {
		fmt.Printf("Invalid birthday format: %v\n", err)
		return
	}

	fmt.Print("Enter role IDs (comma-separated): ")
	roleIDsStr, _ := reader.ReadString('\n')
	roleIDsStr = strings.TrimSpace(roleIDsStr)
	var roleIDs []int64
	if roleIDsStr != "" {
		roleIDParts := strings.Split(roleIDsStr, ",")
		for _, part := range roleIDParts {
			part = strings.TrimSpace(part)
			if roleID, err := strconv.ParseInt(part, 10, 64); err == nil {
				roleIDs = append(roleIDs, roleID)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert roleIDs to strings
	var roleIDStrings []string
	for _, roleID := range roleIDs {
		roleIDStrings = append(roleIDStrings, strconv.FormatInt(roleID, 10))
	}

	resp, err := c.userClient.UpdateUser(ctx, &proto_user.UpdateUserRequest{
		Id:       strconv.FormatInt(userID, 10),
		Email:    email,
		Phone:    phone,
		FullName: fullName,
		Avatar:   avatar,
		Bio:      bio,
		Address:  address,
		Status:   status,
		Birthday: timestamppb.New(birthday),
		RoleIds:  roleIDStrings,
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

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid user ID: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.sessionClient.GetSessionsByUserId(ctx, &proto_session.GetSessionsByUserIdRequest{
		UserId: strconv.FormatInt(userID, 10),
	})
	if err != nil {
		fmt.Printf("Error calling GetSessionsByUserId: %v\n", err)
		return
	}

	fmt.Printf("Found %d sessions for user %d:\n", len(resp.Sessions), userID)
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

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid user ID: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.sessionClient.DeleteSessionByTypeAndUser(ctx, &proto_session.DeleteSessionByTypeAndUserRequest{
		Type:   sessionType,
		UserId: strconv.FormatInt(userID, 10),
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

// Role Service Tests
func (c *GRPCClient) TestGetAllRoles() {
	fmt.Println("\n=== Test GetAllRoles ===")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.roleClient.GetAllRoles(ctx, &proto_role.GetAllRolesRequest{})
	if err != nil {
		fmt.Printf("Error calling GetAllRoles: %v\n", err)
		return
	}

	fmt.Printf("Found %d roles:\n", len(resp.Roles))
	for i, role := range resp.Roles {
		fmt.Printf("Role %d:\n", i+1)
		fmt.Printf("  ID: %s\n", role.Id)
		fmt.Printf("  Name: %s\n", role.Name)
		fmt.Println()
	}
}

func (c *GRPCClient) TestGetRoleById() {
	fmt.Println("\n=== Test GetRoleById ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter role ID: ")
	roleIDStr, _ := reader.ReadString('\n')
	roleIDStr = strings.TrimSpace(roleIDStr)

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid role ID: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.roleClient.GetRoleById(ctx, &proto_role.GetRoleByIdRequest{
		Id: strconv.FormatInt(roleID, 10),
	})
	if err != nil {
		fmt.Printf("Error calling GetRoleById: %v\n", err)
		return
	}

	fmt.Printf("Role found:\n")
	fmt.Printf("ID: %s\n", resp.Role.Id)
	fmt.Printf("Name: %s\n", resp.Role.Name)
}

func (c *GRPCClient) TestCreateRole() {
	fmt.Println("\n=== Test CreateRole ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter role name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter role description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Enter role variant: ")
	variant, _ := reader.ReadString('\n')
	variant = strings.TrimSpace(variant)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.roleClient.CreateRole(ctx, &proto_role.CreateRoleRequest{
		Name:        name,
		Description: description,
		Variant:     variant,
	})
	if err != nil {
		fmt.Printf("Error calling CreateRole: %v\n", err)
		return
	}

	fmt.Printf("Create result: %s\n", resp.Message)
	fmt.Printf("Success: %t\n", resp.Success)
}

func (c *GRPCClient) TestUpdateRole() {
	fmt.Println("\n=== Test UpdateRole ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter role ID to update: ")
	roleIDStr, _ := reader.ReadString('\n')
	roleIDStr = strings.TrimSpace(roleIDStr)

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid role ID: %v\n", err)
		return
	}

	fmt.Print("Enter new role name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter new role description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Enter new role variant: ")
	variant, _ := reader.ReadString('\n')
	variant = strings.TrimSpace(variant)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.roleClient.UpdateRole(ctx, &proto_role.UpdateRoleRequest{
		Id:          strconv.FormatInt(roleID, 10),
		Name:        name,
		Description: description,
		Variant:     variant,
	})
	if err != nil {
		fmt.Printf("Error calling UpdateRole: %v\n", err)
		return
	}

	fmt.Printf("Role updated successfully:\n")
	fmt.Printf("ID: %s\n", resp.Role.Id)
	fmt.Printf("Name: %s\n", resp.Role.Name)
}

func (c *GRPCClient) TestDeleteRole() {
	fmt.Println("\n=== Test DeleteRole ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter role ID to delete: ")
	roleIDStr, _ := reader.ReadString('\n')
	roleIDStr = strings.TrimSpace(roleIDStr)

	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid role ID: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.roleClient.DeleteRole(ctx, &proto_role.DeleteRoleRequest{
		Id: strconv.FormatInt(roleID, 10),
	})
	if err != nil {
		fmt.Printf("Error calling DeleteRole: %v\n", err)
		return
	}

	fmt.Printf("Delete result: %s\n", resp.Message)
	fmt.Printf("Success: %t\n", resp.Success)
}

func printMenu() {
	fmt.Println("\n=== gRPC User Service Test Client ===")
	fmt.Println("1. User Service Tests")
	fmt.Println("  1.1 Get User By ID")
	fmt.Println("  1.2 Delete User By ID")
	fmt.Println("  1.3 Update User By ID")
	fmt.Println("2. Session Service Tests")
	fmt.Println("  2.1 Get All Sessions")
	fmt.Println("  2.2 Get Sessions By User ID")
	fmt.Println("  2.3 Delete Session By Type And Token")
	fmt.Println("  2.4 Delete Session By Type And User")
	fmt.Println("  2.5 Delete Session Expired")
	fmt.Println("3. Role Service Tests")
	fmt.Println("  3.1 Get All Roles")
	fmt.Println("  3.2 Get Role By ID")
	fmt.Println("  3.3 Create Role")
	fmt.Println("  3.4 Update Role")
	fmt.Println("  3.5 Delete Role")
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
		case "1.2":
			client.TestDeleteUserById()
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
		case "3.1":
			client.TestGetAllRoles()
		case "3.2":
			client.TestGetRoleById()
		case "3.3":
			client.TestCreateRole()
		case "3.4":
			client.TestUpdateRole()
		case "3.5":
			client.TestDeleteRole()
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
