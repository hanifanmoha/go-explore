package booking_service

import (
	"context"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/clients"
	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
	repositories_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/repositories"
	services_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/services"
)

type BookingService struct {
	bookingRepo repositories_model.BookingRepository
	movieRepo   repositories_model.MovieRepository
	seatRepo    repositories_model.SeatRepository
	dbClient    *clients.DatabaseClient
}

func NewBookingService(
	bookingRepo repositories_model.BookingRepository,
	movieRepo repositories_model.MovieRepository,
	seatRepo repositories_model.SeatRepository,
	dbClient *clients.DatabaseClient,
) services_model.BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		movieRepo:   movieRepo,
		seatRepo:    seatRepo,
		dbClient:    dbClient,
	}
}

func (s *BookingService) GetBookingsByUserID(ctx context.Context, userID string) ([]entities.Booking, error) {
	return s.bookingRepo.GetBookingsByUserID(ctx, userID)
}

func (s *BookingService) BookSeat(ctx context.Context, movieID int, seatID int, userID string) error {
	return nil
}

func (s *BookingService) OrderSeat(ctx context.Context, movieID int, seatID int, userID string) error {

	tx, err := s.dbClient.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	booking := &entities.Booking{
		MovieID: movieID,
		SeatID:  seatID,
		UserID:  userID,
		Status:  entities.BOOKING_STATUS_PENDING,
	}

	_, err = s.bookingRepo.CreateBooking(ctx, tx, booking)
	if err != nil {
		return err
	}

	err = s.seatRepo.UpdateSeatStatus(ctx, tx, seatID, entities.SEAT_STATUS_SOLD)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
