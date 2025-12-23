package auth_gateway

import (
	"context"
	"net/http"
	"runtime/debug"
)

type basicAuthHandler func(*http.Request) (string, error)

func WithBasicAuth(h http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "No authorization header provided", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		cred, err := Base64Decode(parts[1])
		if err != nil {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		un, pw := strings.SplitN(string(cred), ":", 2)
		if un != username || pw != password {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func recoverPanic(w http.ResponseWriter, err interface{}) {
	if r := recover(); r != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		debug.PrintStack()
	}
}