package grpcservice

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		logRequest(info.FullMethod, req)
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		logResponse(info.FullMethod, resp, err, duration)
		return resp, err
	}
}

func logRequest(method string, req interface{}) {
	maskedReq := maskSensitiveFields(req)
	fmt.Printf("[gRPC] Request - Method: %s, Data: %+v\n", method, maskedReq)
}

func logResponse(method string, resp interface{}, err error, duration time.Duration) {
	statusCode := codes.OK
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		} else {
			statusCode = codes.Internal
		}
	}

	maskedResp := maskSensitiveFields(resp)

	fmt.Printf("[gRPC] Response - Method: %s, Status: %s, Duration: %v, Data: %+v\n",
		method, statusCode, duration, maskedResp)
}

func maskSensitiveFields(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	return maskSensitiveData(data)
}

func maskSensitiveData(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	return "[MASKED]"
}
