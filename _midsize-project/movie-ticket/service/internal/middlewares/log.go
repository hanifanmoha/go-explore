package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		message := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		slog.Info(message)

		next.ServeHTTP(w, r)
	})
}
