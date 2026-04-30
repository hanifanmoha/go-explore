package booking_repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/clients"
	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
	repositories_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/repositories"
	"github.com/jackc/pgx/v5"
)

type BookingRepository struct {
	dbClient *clients.DatabaseClient
}

func NewBookingRepository(dbClient *clients.DatabaseClient) repositories_model.BookingRepository {
	return &BookingRepository{
		dbClient: dbClient,
	}
}

func (r *BookingRepository) GetBookingsByUserID(ctx context.Context, userID string) ([]entities.Booking, error) {

	query := sq.Select("id", "movie_id", "seat_id", "user_id", "status").
		From("bookings").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbClient.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []entities.Booking
	for rows.Next() {
		var booking entities.Booking
		err := rows.Scan(&booking.ID, &booking.MovieID, &booking.SeatID, &booking.UserID, &booking.Status)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepository) CreateBooking(ctx context.Context, tx pgx.Tx, booking *entities.Booking) (*entities.Booking, error) {

	query := sq.Insert("bookings").
		Columns("movie_id", "seat_id", "user_id", "status").
		Values(booking.MovieID, booking.SeatID, booking.UserID, booking.Status).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var bookingID int
	err = tx.QueryRow(ctx, sql, args...).Scan(&bookingID)
	if err != nil {
		return nil, err
	}

	booking.ID = bookingID

	return booking, nil
}
