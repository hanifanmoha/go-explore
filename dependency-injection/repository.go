package dependencyinjection

import "errors"

type Repository interface {
	GetUserByID(userID string) (*User, error)
	UpdateBalanceByUserID(userID string, balance int) error
}

var _ Repository = (*PostgreRepository)(nil)

type PostgreRepository struct {
	users map[string]User
}

func NewPostgreRepository() *PostgreRepository {
	return &PostgreRepository{
		users: map[string]User{
			"id-alice": {
				UserID:  "1",
				Name:    "Alice",
				Balance: 13,
			},
			"id-bob": {
				UserID:  "2",
				Name:    "Bob",
				Balance: 7,
			},
		},
	}
}

func (r *PostgreRepository) GetUserByID(userID string) (*User, error) {
	user, exist := r.users[userID]
	if !exist {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *PostgreRepository) UpdateBalanceByUserID(userID string, balance int) error {
	user, exist := r.users[userID]
	if !exist {
		return errors.New("user not found")
	}
	user.Balance = balance
	r.users[userID] = user
	return nil
}
