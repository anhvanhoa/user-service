package grpcservice

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor provides logging for gRPC services
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		
		// Log request
		logRequest(info.FullMethod, req)
		
		// Call the next handler
		resp, err := handler(ctx, req)
		
		// Calculate duration
		duration := time.Since(start)
		
		// Log response
		logResponse(info.FullMethod, resp, err, duration)
		
		return resp, err
	}
}

// logRequest logs the incoming request
func logRequest(method string, req interface{}) {
	// Mask sensitive fields
	maskedReq := maskSensitiveFields(req)
	
	fmt.Printf("[gRPC] Request - Method: %s, Data: %+v\n", method, maskedReq)
}

// logResponse logs the response
func logResponse(method string, resp interface{}, err error, duration time.Duration) {
	statusCode := codes.OK
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		} else {
			statusCode = codes.Internal
		}
	}
	
	// Mask sensitive fields in response
	maskedResp := maskSensitiveFields(resp)
	
	fmt.Printf("[gRPC] Response - Method: %s, Status: %s, Duration: %v, Data: %+v\n", 
		method, statusCode, duration, maskedResp)
}

// maskSensitiveFields masks sensitive information in requests/responses
func maskSensitiveFields(data interface{}) interface{} {
	if data == nil {
		return nil
	}
	
	// Use reflection to mask sensitive fields
	// This is a simplified version - in production, you might want to use a more sophisticated approach
	return maskSensitiveData(data)
}

// maskSensitiveData recursively masks sensitive fields
func maskSensitiveData(data interface{}) interface{} {
	if data == nil {
		return nil
	}
	
	// For now, we'll just return a generic masked response
	// In a real implementation, you would use reflection to identify and mask specific fields
	return "[MASKED]"
}

// Structured logging with different levels
func logInfo(method, message string, fields map[string]interface{}) {
	fmt.Printf("[INFO] %s - %s: %+v\n", method, message, fields)
}

func logError(method, message string, err error, fields map[string]interface{}) {
	fmt.Printf("[ERROR] %s - %s: %v, Fields: %+v\n", method, message, err, fields)
}

func logWarn(method, message string, fields map[string]interface{}) {
	fmt.Printf("[WARN] %s - %s: %+v\n", method, message, fields)
}

func logDebug(method, message string, fields map[string]interface{}) {
	fmt.Printf("[DEBUG] %s - %s: %+v\n", method, message, fields)
}
