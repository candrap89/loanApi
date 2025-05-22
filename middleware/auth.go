package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	redis "github.com/candrap89/loanApi/Redis"
	"github.com/candrap89/loanApi/queries"
)

type RouteAPIKey struct {
	Route     string
	APIKey    string
	RateLimit int
}
type contextKey string

const UserLoanKey contextKey = "user_loan_data"

// APIKeyMiddleware verifies the presence and validity of the API key
func APIKeyMiddleware(routeKeys []RouteAPIKey, redisClient *redis.RedisClient, userLoanQuery queries.UserLoanQueryInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// use context to pass data between middleware
			ctx := r.Context()
			route := r.URL.Path

			cif := r.Header.Get("CIF")
			if cif != "" {
				userLoan, err := userLoanQuery.GetUserLoanByCIF(cif)
				if err != nil {
					fmt.Println("user loan  ", userLoan)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			// Find matching route config
			var routeConfig *RouteAPIKey
			for _, rk := range routeKeys {
				if rk.Route == route {
					routeConfig = &rk
					fmt.Println("route apikey config ", routeConfig.APIKey)
					break
				}
			}

			if routeConfig == nil {
				http.Error(w, "Route not configured", http.StatusNotFound)
				return
			}

			// Verify API Key
			apiKey := r.Header.Get("X-API-Key")
			if apiKey != routeConfig.APIKey {
				fmt.Println("api key ", apiKey)
				fmt.Println("api key from config ", routeConfig.APIKey)
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// Rate Limiting
			redisKey := "rate_limit:" + apiKey + ":" + route
			count, err := redisClient.IncrementCounter(ctx, redisKey, 60*time.Second)
			if err != nil {
				fmt.Println("Error incrementing counter: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if count > int64(routeConfig.RateLimit) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			//Add hit count to context
			ctx = context.WithValue(ctx, "hit_count", count)
			r = r.WithContext(ctx)
			fmt.Println("Hit count: ", ctx.Value("hit_count"))

			next.ServeHTTP(w, r)
		})
	}
}

func GetLoanDataFromContext(w http.ResponseWriter, r *http.Request) {
	loanData := r.Context().Value(UserLoanKey)
	if loanData == nil {
		http.Error(w, "No user loan data in context", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loanData)
}
