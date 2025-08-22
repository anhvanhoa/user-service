package grpcservice

import (
	"context"
	"fmt"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationInterceptor provides automatic validation for gRPC requests
func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Validate the request if it has a Validate method
		if err := validateRequest(req); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
		}

		// Call the next handler
		return handler(ctx, req)
	}
}

// validateRequest validates a request using its Validate method if available
func validateRequest(req interface{}) error {
	if req == nil {
		return nil
	}

	// Use reflection to check if the request has a Validate method
	reqType := reflect.TypeOf(req)
	reqValue := reflect.ValueOf(req)

	// Check if the request is a pointer
	if reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem()
		reqValue = reqValue.Elem()
	}

	// Look for Validate method
	validateMethod := reqValue.MethodByName("Validate")
	if !validateMethod.IsValid() {
		// If no Validate method, try ValidateAll
		validateMethod = reqValue.MethodByName("ValidateAll")
		if !validateMethod.IsValid() {
			// No validation method found, skip validation
			return nil
		}
	}

	// Call the validation method
	results := validateMethod.Call(nil)
	if len(results) > 0 {
		if err := results[0].Interface(); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	return nil
}
