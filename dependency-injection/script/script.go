package main

import (
	"fmt"

	di "github.com/hanifanmoha/go-explore/dependency-injection"
)

const (
	ID_ALICE = "id-alice"
	ID_BOB   = "id-bob"
)

func main() {
	pgRepo := di.NewPostgreRepository()
	svc := di.Service{
		Repo: pgRepo,
	}

	// Test GetBalance
	alice, err := svc.GetBalance(ID_ALICE)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Initial Alice: %+v\n", alice)

	// Test Deposit
	err = svc.Deposit(ID_ALICE, 5)
	if err != nil {
		panic(err)
	}
	alice, _ = svc.GetBalance(ID_ALICE)
	fmt.Printf("After Deposit 5 to Alice: %+v\n", alice)

	// Test Withdraw
	err = svc.Withdraw(ID_ALICE, 7)
	if err != nil {
		panic(err)
	}
	alice, _ = svc.GetBalance(ID_ALICE)
	fmt.Printf("After Withdraw 7 from Alice: %+v\n", alice)

	// Test Transfer
	err = svc.Transfer(ID_ALICE, ID_BOB, 3)
	if err != nil {
		panic(err)
	}
	alice, _ = svc.GetBalance(ID_ALICE)
	bob, _ := svc.GetBalance(ID_BOB)
	fmt.Printf("After Transfer 3 from Alice to Bob:\nAlice: %+v\nBob: %+v\n", alice, bob)
}
