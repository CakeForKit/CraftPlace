package loggergrpc

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// GRPCLoggingInterceptor создает интерцептор для логирования gRPC запросов с logrus
func GRPCLoggingInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Получаем информацию о клиенте
		clientIP := getGRPCClientIP(ctx)

		// Создаем логгер для этого запроса
		requestLogger := logger.WithFields(logrus.Fields{
			"type":      "grpc",
			"method":    info.FullMethod,
			"client_ip": clientIP,
		})

		// Логируем начало запроса (только для не-чувствительных методов)
		if !isSensitiveMethod(info.FullMethod) {
			requestLogger.WithFields(logrus.Fields{
				"request": maskSensitiveData(info.FullMethod, req),
			}).Info("gRPC request started")
		} else {
			requestLogger.Info("gRPC request started")
		}

		// Обрабатываем запрос
		resp, err := handler(ctx, req)

		// Логируем результат
		duration := time.Since(start)
		grpcStatus := status.Code(err)

		logFields := logrus.Fields{
			"duration":       duration.Milliseconds(),
			"duration_human": duration.String(),
			"grpc_status":    grpcStatus.String(),
			"grpc_code":      grpcStatus,
		}

		// Добавляем информацию об ответе для не-чувствительных методов
		if !isSensitiveMethod(info.FullMethod) && err == nil {
			logFields["response"] = maskSensitiveData(info.FullMethod, resp)
		}

		if err != nil {
			requestLogger.WithFields(logFields).WithError(err).Error("gRPC request failed")
		} else {
			requestLogger.WithFields(logFields).Info("gRPC request completed")
		}

		return resp, err
	}
}

// GRPCStreamLoggingInterceptor для streaming запросов
func GRPCStreamLoggingInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		clientIP := getGRPCClientIP(ss.Context())

		requestLogger := logger.WithFields(logrus.Fields{
			"type":             "grpc_stream",
			"method":           info.FullMethod,
			"client_ip":        clientIP,
			"is_client_stream": info.IsClientStream,
			"is_server_stream": info.IsServerStream,
		})

		requestLogger.Info("gRPC stream started")

		err := handler(srv, ss)

		duration := time.Since(start)
		grpcStatus := status.Code(err)

		logFields := logrus.Fields{
			"duration":       duration.Milliseconds(),
			"duration_human": duration.String(),
			"grpc_status":    grpcStatus.String(),
		}

		if err != nil {
			requestLogger.WithFields(logFields).WithError(err).Error("gRPC stream failed")
		} else {
			requestLogger.WithFields(logFields).Info("gRPC stream completed")
		}

		return err
	}
}

// getGRPCClientIP извлекает IP клиента из контекста gRPC
func getGRPCClientIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}
	return "unknown"
}

// isSensitiveMethod проверяет, является ли метод чувствительным (не логируем тело)
func isSensitiveMethod(method string) bool {
	sensitiveMethods := map[string]bool{
		"/auth.AuthService/LoginUser":    true,
		"/auth.AuthService/RegisterUser": true,
		// Добавьте другие чувствительные методы
	}
	return sensitiveMethods[method]
}

// maskSensitiveData маскирует чувствительные данные в логах
func maskSensitiveData(method string, data interface{}) interface{} {
	// Для чувствительных методов возвращаем маскированную версию
	if isSensitiveMethod(method) {
		return "[SENSITIVE_DATA_MASKED]"
	}
	return data
}
