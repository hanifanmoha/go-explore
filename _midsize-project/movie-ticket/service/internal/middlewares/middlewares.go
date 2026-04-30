package middlewares

import "net/http"

func Middlewares(next http.Handler) http.Handler {
	next = logMiddleware(next)
	next = corsMiddleware(next)
	return next
}
