package main

import (
	"fmt"
	"net/http"
	"time"
)

func SendEmail(to string, subject string) {
	fmt.Printf("Sending email to %s with subject '%s ...'\n", to, subject)
	// Simulate email sending delay
	time.Sleep(3 * time.Second)
	fmt.Println("Email sent!")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	fmt.Println("Record user data to database.")
	fmt.Println("Create other user related data.")
	fmt.Println("Do other things.")

	// Add task to queue instead of calling directly
	// SendEmail("user@example.com", "Welcome to async party!")
	AddTaskEmailSend("user@example.com", "Welcome to async party!")

	elapsed := time.Since(start)
	fmt.Fprintf(w, "User Created! Time taken: %s", elapsed)
}

func main() {
	StartWorker()
	http.HandleFunc("/register", RegisterHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
