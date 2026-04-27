package chainofresponsibility

import (
	"fmt"
	"slices"
)

// https://refactoring.guru/design-patterns/chain-of-responsibility/go/example

// ====== MIDDLEWARE INTERFACE =====

// implement the chain of responsibility pattern for handling requests
type Middleware interface {
	execute(request *Request) error
	setNext(m Middleware)
}

// ====== AUTH MIDDLEWARE =====

type AuthMiddleware struct {
	next               Middleware
	token              string
	authenticatedPaths []string
}

func NewAuthMiddleware(token string, authenticatedPaths []string) *AuthMiddleware {
	return &AuthMiddleware{
		token:              "super-secret-token",
		authenticatedPaths: []string{"/profiles", "/users"},
	}
}

func (a *AuthMiddleware) execute(request *Request) error {
	isAuthenticatedPath := slices.Contains(a.authenticatedPaths, request.Path)
	if !isAuthenticatedPath {
		return nil
	}
	if request.Token != a.token {
		return fmt.Errorf("unauthorized")
	}

	// continue to the next middleware in the chain
	if a.next != nil {
		return a.next.execute(request)
	}
	return nil
}

func (a *AuthMiddleware) setNext(m Middleware) {
	a.next = m
}

// ====== ROLE MIDDLEWARE =====

type RoleMiddleware struct {
	next Middleware
}

func NewRoleMiddleware() *RoleMiddleware {
	return &RoleMiddleware{}
}

func (r *RoleMiddleware) execute(request *Request) error {
	if request.Path == "/users" && request.Role != "admin" {
		return fmt.Errorf("forbidden")
	}

	// continue to the next middleware in the chain
	if r.next != nil {
		return r.next.execute(request)
	}
	return nil
}

func (r *RoleMiddleware) setNext(m Middleware) {
	r.next = m
}

// ====== LOGGER MIDDLEWARE =====

type LoggerMiddleware struct {
	next Middleware
}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (l *LoggerMiddleware) execute(request *Request) error {
	fmt.Printf("Request path: %s, role: %s\n", request.Path, request.Role)

	// continue to the next middleware in the chain
	if l.next != nil {
		return l.next.execute(request)
	}
	return nil
}

func (l *LoggerMiddleware) setNext(m Middleware) {
	l.next = m
}

// ===== REQUEST STRUCT =====

type Request struct {
	Path  string
	Role  string
	Token string
}

// ====== HANDLER =====

type Handler struct {
	middleware Middleware
}

func NewHandler() *Handler {

	roleMiddleware := NewRoleMiddleware()
	authMiddleware := NewAuthMiddleware("super-secret-token", []string{"/profiles", "/users"})
	authMiddleware.setNext(roleMiddleware)
	loggerMiddleware := NewLoggerMiddleware()
	loggerMiddleware.setNext(authMiddleware)

	return &Handler{
		middleware: loggerMiddleware,
	}
}

func (h *Handler) Execute(request *Request) error {
	if h.middleware != nil {
		return h.middleware.execute(request)
	}
	return nil
}
