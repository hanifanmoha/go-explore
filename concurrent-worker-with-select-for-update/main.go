package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Worker struct {
	ID     uint   `gorm:"primaryKey"`
	Status string `gorm:"not null;default:pending"`
}

type WorkerLog struct {
	ID         uint       `gorm:"primaryKey"`
	WorkerID   uint       `gorm:"not null"`
	FinishedAt *time.Time `gorm:"type:timestamptz"`
	WorkerName string     `gorm:"not null"`
}

func main() {
	// Database connection
	db, err := setupDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	ctx := context.Background()

	shutdown := make(chan int)
	go run(ctx, shutdown, db)

	log.Println("Worker is running...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	log.Println("Press Ctrl+C to stop the worker")
	<-quit

	shutdown <- 1

	log.Println("Worker is shutting down")
}

func setupDatabase() (*gorm.DB, error) {
	// Get database connection info from environment variables
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		"postgres", "worker", "password", "workerdb", "5432")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

func run(ctx context.Context, shutdown chan int, db *gorm.DB) {
	for {
		select {
		case <-shutdown:
			log.Println("Shutting down...")
			return
		default:
			log.Println("Run2!")
			runWorker(ctx, db)
			// time.Sleep(5 * time.Second)
		}
	}
}

func runWorker(ctx context.Context, db *gorm.DB) {
	workerName := getEnv("WORKER_ID", "unknown-worker")
	var worker Worker

	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// With row locking, worker logs will be unique
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).Where("status = ?", "pending").First(&worker).Error; err != nil {
			return err
		}
		// Without row locking, worker logs will be duplicated
		// if err := tx.Where("status = ?", "pending").First(&worker).Error; err != nil {
		// 	return err
		// }

		log.Printf("%s : Acquired lock on worker ID: %d\n", workerName, worker.ID)

		worker.Status = "finished"

		if err := tx.Save(&worker).Error; err != nil {
			return err
		}

		// save worker logs
		now := time.Now()
		WorkerLog := WorkerLog{
			WorkerID:   worker.ID,
			FinishedAt: &now,
			WorkerName: workerName,
		}

		if err := tx.Create(&WorkerLog).Error; err != nil {
			return err
		}

		// random sleep between 2-5 seconds to simulate work
		sleepDuration := time.Duration(2+time.Now().UnixNano()%10) * time.Second
		time.Sleep(sleepDuration)

		return nil
	}); err != nil {
		log.Println("No pending worker found or error occurred:", err)
		return
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
