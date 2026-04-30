package repositories_model

import (
	"context"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
	"github.com/jackc/pgx/v5"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, tx pgx.Tx, booking *entities.Booking) (*entities.Booking, error)
}
