package dependencyinjection

import "errors"

type Service struct {
	Repo Repository
}

func (s *Service) GetBalance(userID string) (*User, error) {
	return s.Repo.GetUserByID(userID)
}

func (s *Service) Deposit(userID string, val int) error {

	if val < 0 {
		return errors.New("invalid value")
	}

	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	err = s.Repo.UpdateBalanceByUserID(userID, user.Balance+val)
	return err
}

func (s *Service) Withdraw(userID string, val int) error {

	if val < 0 {
		return errors.New("invalid value")
	}

	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.Balance < val {
		return errors.New("insufficient balance")
	}

	return s.Repo.UpdateBalanceByUserID(userID, user.Balance-val)
}

func (s *Service) Transfer(senderID string, receiverID string, val int) error {

	if val < 0 {
		return errors.New("invalid value")
	}

	err := s.Withdraw(senderID, val)
	if err != nil {
		return err
	}
	err = s.Deposit(receiverID, val)
	if err != nil {
		return err
	}

	return nil
}
