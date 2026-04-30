package services_model

import (
	"context"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
)

type BookingService interface {
	GetBookingsByUserID(ctx context.Context, userID string) ([]entities.Booking, error)
	BookSeat(ctx context.Context, movieID int, seatID int, userID string) error
	OrderSeat(ctx context.Context, movieID int, seatID int, userID string) error
}
