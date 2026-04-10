package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// == MODEL ==================================================

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	OrganizationID string        `json:"organization_id"`
	Organization   *Organization `json:"organization,omitempty"`
}

// == REPOSITORY ==============================================

type UserRepository struct {
	users map[int]*User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: map[int]*User{
			1: {ID: 1, Name: "Agus", OrganizationID: "org1"},
			2: {ID: 2, Name: "Beni", OrganizationID: "org2"},
			3: {ID: 3, Name: "Caca", OrganizationID: "org3"},
		},
	}
}

func (r *UserRepository) GetUserByID(id int) (*User, error) {
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

type OrganizationRepository struct {
	organizations map[string]*Organization
}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{
		organizations: map[string]*Organization{
			"org1": {ID: "org1", Name: "Bersahaja Corp"},
			"org2": {ID: "org2", Name: "PT Tumbuh Bestari"},
		},
	}
}

func (r *OrganizationRepository) GetOrganizationByID(id string) (*Organization, error) {
	if org, ok := r.organizations[id]; ok {
		return org, nil
	}
	return nil, fmt.Errorf("organization not found")
}

// == SERVICE =================================================

type UserService struct {
	userRepo         *UserRepository
	organizationRepo *OrganizationRepository
}

func NewUserService(userRepo *UserRepository, organizationRepo *OrganizationRepository) *UserService {
	return &UserService{
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
	}
}

func (s *UserService) GetUserWithOrganization(id int) (*User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	org, err := s.organizationRepo.GetOrganizationByID(user.OrganizationID)
	if err != nil {
		return nil, err
	}
	user.Organization = org
	return user, nil
}

// == HANDLER =================================================

type Handler struct {
	userService *UserService
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserWithOrganization(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// == MAIN ====================================================

func main() {
	userRepo := NewUserRepository()
	orgRepo := NewOrganizationRepository()
	userService := NewUserService(userRepo, orgRepo)
	handler := &Handler{userService: userService}

	http.HandleFunc("/user/{id}", handler.GetUser)

	http.ListenAndServe(":8080", nil)
}
