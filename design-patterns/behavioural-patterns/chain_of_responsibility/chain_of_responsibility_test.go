package chainofresponsibility_test

import (
	"testing"

	chainofresponsibility "github.com/hanifanmoha/go-explore/design-patterns/behavioural-patterns/chain_of_responsibility"
)

func TestHandler_Execute(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		request *chainofresponsibility.Request
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := chainofresponsibility.NewHandler()
			gotErr := h.Execute(tt.request)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Execute() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Execute() succeeded unexpectedly")
			}
		})
	}

	t.Run("Unauthorized request cannot access /profiles", func(t *testing.T) {
		request := &chainofresponsibility.Request{
			Path:  "/profiles",
			Token: "invalid-token",
		}

		h := chainofresponsibility.NewHandler()
		err := h.Execute(request)
		if err == nil {
			t.Fatal("Expected error for unauthorized request, got nil")
		}
	})

	t.Run("Unauthorized request cannot access /users", func(t *testing.T) {
		request := &chainofresponsibility.Request{
			Path:  "/users",
			Token: "invalid-token",
		}

		h := chainofresponsibility.NewHandler()
		err := h.Execute(request)
		if err == nil {
			t.Fatal("Expected error for unauthorized request, got nil")
		}
	})

	t.Run("Authorized request can access /profiles", func(t *testing.T) {
		request := &chainofresponsibility.Request{
			Path:  "/profiles",
			Token: "super-secret-token",
		}

		h := chainofresponsibility.NewHandler()
		err := h.Execute(request)
		if err != nil {
			t.Fatalf("Expected no error for authorized request, got %v", err)
		}
	})

	t.Run("Authorized request can access /users", func(t *testing.T) {
		request := &chainofresponsibility.Request{
			Path:  "/users",
			Token: "super-secret-token",
			Role:  "admin",
		}

		h := chainofresponsibility.NewHandler()
		err := h.Execute(request)
		if err != nil {
			t.Fatalf("Expected no error for authorized request, got %v", err)
		}
	})

	t.Run("Non-admin user cannot access /users", func(t *testing.T) {
		request := &chainofresponsibility.Request{
			Path:  "/users",
			Token: "super-secret-token",
			Role:  "user",
		}

		h := chainofresponsibility.NewHandler()
		err := h.Execute(request)
		if err == nil {
			t.Fatal("Expected error for non-admin user, got nil")
		}
	})
}
