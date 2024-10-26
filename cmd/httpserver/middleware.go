package httpserver

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go_chi_pgx/state"
	utils "go_chi_pgx/utils"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

func AuthMiddleware(app *state.State) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractTokenFromHeader(r)
			if tokenStr == "" {
				app.Logger.PrintError(fmt.Errorf("no token provided"), map[string]string{
					"context": "authorization",
				})
				_ = Unauthorized.WriteToResponse(w, nil)
				return
			}

			var claims utils.Claims
			token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.Config.SecretKey), nil
			})

			if err != nil || !token.Valid {
				app.Logger.PrintError(fmt.Errorf("invalid token"), map[string]string{
					"context": "authorization",
				})
				_ = Unauthorized.WriteToResponse(w, nil)
				return
			}
			ctx := context.WithValue(r.Context(), "userid", claims.UserID.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RateLimitMiddleware(app *state.State) func(http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Goroutine to periodically clean up stale entries
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if app.Config.LimiterEnabled {
				ip, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					app.Logger.PrintError(err, map[string]string{
						"context": "host port split error",
					})
					_ = Unauthorized.WriteToResponse(w, nil)
					return
				}

				mu.Lock()
				cli, found := clients[ip]
				if !found {
					cli = &client{
						limiter: rate.NewLimiter(
							rate.Limit(app.Config.Rps),
							app.Config.Burst,
						),
					}
					clients[ip] = cli
				}
				cli.lastSeen = time.Now()

				if !cli.limiter.Allow() {
					mu.Unlock()
					_ = RateLimitExceeded.WriteToResponse(w, nil)
					return
				}
				mu.Unlock()
			}

			// Proceed with the request
			next.ServeHTTP(w, r)
		})
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userid").(string)
	return userID, ok
}
