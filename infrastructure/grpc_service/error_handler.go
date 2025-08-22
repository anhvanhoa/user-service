package grpcservice

import (
	"context"
	"fmt"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorHandlerInterceptor provides error handling for gRPC services
func ErrorHandlerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic
				fmt.Printf("Panic in %s: %v\n%s", info.FullMethod, r, debug.Stack())
			}
		}()

		// Call the next handler
		resp, err := handler(ctx, req)
		if err != nil {
			// Convert domain errors to gRPC status errors
			return nil, convertError(err)
		}

		return resp, nil
	}
}

// convertError converts domain errors to gRPC status errors
func convertError(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's already a gRPC status error
	if _, ok := status.FromError(err); ok {
		return err
	}

	// Convert common domain errors
	errMsg := err.Error()
	
	// Authentication errors
	if contains(errMsg, "không tìm thấy") || contains(errMsg, "not found") {
		return status.Errorf(codes.NotFound, errMsg)
	}
	
	// Validation errors
	if contains(errMsg, "không hợp lệ") || contains(errMsg, "invalid") || contains(errMsg, "không đúng") {
		return status.Errorf(codes.InvalidArgument, errMsg)
	}
	
	// Authorization errors
	if contains(errMsg, "không có quyền") || contains(errMsg, "unauthorized") {
		return status.Errorf(codes.PermissionDenied, errMsg)
	}
	
	// Already exists errors
	if contains(errMsg, "đã tồn tại") || contains(errMsg, "already exists") {
		return status.Errorf(codes.AlreadyExists, errMsg)
	}
	
	// Resource exhausted errors
	if contains(errMsg, "hết hạn") || contains(errMsg, "expired") {
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}
	
	// Database errors
	if contains(errMsg, "database") || contains(errMsg, "connection") {
		return status.Errorf(codes.Unavailable, "Lỗi hệ thống, vui lòng thử lại sau")
	}
	
	// Default to internal error
	return status.Errorf(codes.Internal, "Lỗi hệ thống: %s", errMsg)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

// containsSubstring checks if a string contains a substring (case-insensitive)
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// BusinessError represents a business logic error
type BusinessError struct {
	Code    string
	Message string
}

func (e BusinessError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewBusinessError creates a new business error
func NewBusinessError(code, message string) error {
	return BusinessError{
		Code:    code,
		Message: message,
	}
}
