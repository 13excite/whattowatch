package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var reqid uint64
var prefix string

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey int = 0

// RequestIDHeader is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var RequestIDHeader = "X-Request-Id"

type (
	// struct for holding response details
	responseData struct {
		status int
		size   int
	}

	// our http.ResponseWriter implementation
	loggingResponseWriter struct {
		http.ResponseWriter // compose original http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.responseData.size += size            // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.responseData.status = statusCode       // capture status code
}

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			myid := atomic.AddUint64(&reqid, 1)
			requestID = fmt.Sprintf("%s-%06d", prefix, myid)
		}
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// NextRequestID generates the next request ID in the sequence.
func NextRequestID() uint64 {
	return atomic.AddUint64(&reqid, 1)
}

func loggerHTTPMiddlewareDefault(disabledEndpoints []string) func(http.Handler) http.Handler {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If Disabled
			if _, ok := disabled[r.RequestURI]; ok {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}

			lrw := loggingResponseWriter{
				ResponseWriter: w, // compose original http.ResponseWriter
				responseData:   responseData,
			}

			next.ServeHTTP(&lrw, r)

			fields := []zapcore.Field{
				zap.Int("status", lrw.responseData.status),
				zap.Duration("duration", time.Since(start)),
				zap.String("path", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("package", "server.http"),
			}

			if reqID := r.Context().Value(RequestIDKey); reqID != nil {
				fields = append(fields, zap.String("request-id", reqID.(string)))
			}

			// If we have an x-Forwarded-For header, use that for the remote
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				fields = append(fields, zap.String("remote", forwardedFor))
			} else {
				fields = append(fields, zap.String("remote", r.RemoteAddr))
			}
			zap.L().Info("HTTP Request", fields...)
		})
	}
}

// Returns a middleware function for logging requests
func loggerHTTPMiddlewareStackdriver(disabledEndpoints []string) func(http.Handler) http.Handler {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If Disabled
			if _, ok := disabled[r.RequestURI]; ok {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}

			lrw := loggingResponseWriter{
				ResponseWriter: w, // compose original http.ResponseWriter
				responseData:   responseData,
			}

			next.ServeHTTP(&lrw, r)

			// If the remote IP is being proxied, use the real IP
			remoteIP := r.Header.Get("X-Forwarded-For")
			if remoteIP == "" {
				remoteIP = r.RemoteAddr
			}

			fields := []zapcore.Field{
				zapdriver.HTTP(&zapdriver.HTTPPayload{
					RequestMethod: r.Method,
					RequestURL:    r.RequestURI,
					RequestSize:   strconv.FormatInt(r.ContentLength, 10),
					Status:        lrw.responseData.status,
					UserAgent:     r.UserAgent(),
					RemoteIP:      remoteIP,
					Referer:       r.Referer(),
					Latency:       fmt.Sprintf("%fs", time.Since(start).Seconds()),
					Protocol:      r.Proto,
				}),
				zap.String("package", "server.http"),
			}

			if reqID := r.Context().Value(RequestIDKey); reqID != nil {
				fields = append(fields, zap.String("request-id", reqID.(string)))
			}

			zap.L().Info("HTTP Request", fields...)
		})
	}
}
