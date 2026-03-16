package main

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

var redisClient = asynq.RedisClientOpt{Addr: "127.0.0.1:6379"}
var asynqClient = asynq.NewClient(redisClient)

const TASK_EMAIL_SEND = "email:send"

type EmailPayload struct {
	To      string
	Subject string
}

func AddTaskEmailSend(to string, subject string) {
	payload, _ := json.Marshal(EmailPayload{
		To:      to,
		Subject: subject,
	})
	task := asynq.NewTask(TASK_EMAIL_SEND, payload)
	asynqClient.Enqueue(task)
}

func HandleEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload EmailPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	SendEmail(payload.To, payload.Subject)
	return nil
}

func StartWorker() {
	go func() {
		srv := asynq.NewServer(
			&redisClient,
			asynq.Config{
				Concurrency: 3, // number of concurrent workers
			},
		)

		mux := asynq.NewServeMux()
		// mux.HandleFunc(TASK_EMAIL_SEND, HandleEmailTask)

		if err := srv.Run(mux); err != nil {
			panic(err)
		}
	}()
}
