package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"Gofinal/pkg/common"

	"github.com/dgrijalva/jwt-go"
)

// ─── Context keys ────────────────────────────────────────────────────────────

// Use context keys from common package
var (
	ContextUserID    = common.ContextUserID
	ContextUserRole  = common.ContextUserRole
	ContextUserEmail = common.ContextUserEmail
	ContextRequestID = common.ContextRequestID
)

// AdminEmail - the only administrator's email
const AdminEmail = "abzalbahktiarow2006@gmail.com"

// JWT secret – must match the one used in auth.Handler.Login
var jwtSecret = []byte("yourSecretKey")

// Middleware is a standard middleware signature.
type Middleware func(http.Handler) http.Handler

// ─── Chain helper ────────────────────────────────────────────────────────────

// Chain composes middlewares left-to-right around a final handler.
//
//	Chain(handler, A, B, C) → A( B( C( handler ) ) )
func Chain(handler http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}

// WrapFunc applies middlewares to a single http.HandlerFunc.
func WrapFunc(fn http.HandlerFunc, mws ...Middleware) http.Handler {
	var h http.Handler = fn
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// ─── Logging ─────────────────────────────────────────────────────────────────

// responseWriter wraps http.ResponseWriter to capture status code and bytes.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

// Logging logs every request: method, remote address, path, status, duration,
// response size.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := newResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		log.Printf("[%s] %s %s | %d | %v | %d bytes",
			r.Method,
			r.RemoteAddr,
			r.URL.Path,
			wrapped.statusCode,
			time.Since(start),
			wrapped.bytes,
		)
	})
}

// ─── Recovery ────────────────────────────────────────────────────────────────

// Recovery catches panics in downstream handlers and returns 500.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %s %s: %v\n%s",
					r.Method, r.URL.Path, err, debug.Stack())
				writeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ─── CORS ────────────────────────────────────────────────────────────────────

// CORS adds permissive CORS headers and handles preflight OPTIONS requests.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ─── Request ID ──────────────────────────────────────────────────────────────

var (
	reqCounter   uint64
	reqCounterMu sync.Mutex
)

// RequestID generates a unique request identifier, puts it in context and in
// the X-Request-ID response header.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCounterMu.Lock()
		reqCounter++
		id := fmt.Sprintf("req-%d-%d", time.Now().UnixNano(), reqCounter)
		reqCounterMu.Unlock()

		ctx := context.WithValue(r.Context(), ContextRequestID, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ─── JWT Authentication ─────────────────────────────────────────────────────

// Authenticate reads the Authorization: Bearer <token> header, validates the
// JWT, and injects userID (int) and userRole (string) into the request context.
// It does NOT reject unauthenticated requests – combine with RequireAuth for
// that.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := parts[1]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		if id, ok := claims["id"].(float64); ok {
			ctx = context.WithValue(ctx, ContextUserID, int(id))
		}
		if role, ok := claims["role"].(string); ok {
			ctx = context.WithValue(ctx, ContextUserRole, role)
		}
		if email, ok := claims["email"].(string); ok {
			ctx = context.WithValue(ctx, ContextUserEmail, email)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth rejects unauthenticated requests with 401.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(ContextUserID) == nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "authentication required",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireRole returns middleware that only allows users whose role matches.
func RequireRole(role string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, _ := r.Context().Value(ContextUserRole).(string)
			if userRole != role {
				writeJSON(w, http.StatusForbidden, map[string]string{
					"error": "forbidden: requires " + role + " role",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin - middleware that checks if user is administrator (email = AdminEmail)
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, _ := r.Context().Value(ContextUserEmail).(string)
		if email != AdminEmail {
			// Redirect to home if not admin
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdminAPI - middleware for API endpoints (returns JSON error)
func RequireAdminAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, _ := r.Context().Value(ContextUserEmail).(string)
		if email != AdminEmail {
			writeJSON(w, http.StatusForbidden, map[string]string{
				"error": "access denied: admin only",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ─── Rate Limiter (per-IP, token bucket) ─────────────────────────────────────

type visitor struct {
	tokens   float64
	lastSeen time.Time
}

// RateLimiterConfig configures the per-IP rate limiter.
type RateLimiterConfig struct {
	RequestsPerSecond float64 // refill rate
	Burst             int     // max tokens (bucket size)
}

// RateLimiter returns middleware that limits requests per IP.
func RateLimiter(cfg RateLimiterConfig) Middleware {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	// Background cleaner – removes stale entries every minute.
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if idx := strings.LastIndex(ip, ":"); idx != -1 {
				ip = ip[:idx]
			}

			mu.Lock()
			v, exists := visitors[ip]
			if !exists {
				v = &visitor{tokens: float64(cfg.Burst)}
				visitors[ip] = v
			}

			elapsed := time.Since(v.lastSeen).Seconds()
			v.lastSeen = time.Now()
			v.tokens += elapsed * cfg.RequestsPerSecond
			if v.tokens > float64(cfg.Burst) {
				v.tokens = float64(cfg.Burst)
			}

			if v.tokens < 1 {
				mu.Unlock()
				writeJSON(w, http.StatusTooManyRequests, map[string]string{
					"error": "rate limit exceeded",
				})
				return
			}

			v.tokens--
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

// ─── Content-Type ────────────────────────────────────────────────────────────

// ContentTypeJSON sets the Content-Type header to application/json for every
// response.
func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

// writeJSON is a small helper used by middleware to return JSON error bodies.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// GetUserIDFromContext extracts the user ID placed by the Authenticate
// middleware.
func GetUserIDFromContext(r *http.Request) (int, bool) {
	id, ok := r.Context().Value(ContextUserID).(int)
	return id, ok
}

// GetUserRoleFromContext extracts the user role placed by the Authenticate
// middleware.
func GetUserRoleFromContext(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(ContextUserRole).(string)
	return role, ok
}

// GetRequestIDFromContext extracts the request ID placed by the RequestID
// middleware.
func GetRequestIDFromContext(r *http.Request) string {
	id, _ := r.Context().Value(ContextRequestID).(string)
	return id
}
