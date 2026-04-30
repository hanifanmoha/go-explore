package services_model

import "context"

type BookingService interface {
	BookSeat(ctx context.Context, movieID int, seatID int, userID string) error
	OrderSeat(ctx context.Context, movieID int, seatID int, userID string) error
}
