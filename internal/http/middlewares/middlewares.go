package middlewares

import "net/http"


func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		// Example: Very basic token check (in production, validate properly!)
		const dummyToken = "mysecrettoken123"
		if authHeader != "Bearer "+dummyToken {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Token valid, proceed to handler
		next.ServeHTTP(w, r)
	}
}