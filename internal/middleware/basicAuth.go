package middleware

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
)

// basicAuthMiddleware protects an http.Handler with Basic Auth (constant-time compare).
func BasicAuthMiddleware(user, pass string, next http.Handler) http.Handler {
	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pass)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if subtle.ConstantTimeCompare([]byte(auth), []byte(expectedAuth)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
