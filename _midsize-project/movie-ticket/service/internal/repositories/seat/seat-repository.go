package seat_repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/clients"
	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
	repositories_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/repositories"
)

type SeatRepository struct {
	dbClient *clients.DatabaseClient
}

func NewSeatRepository(dbClient *clients.DatabaseClient) repositories_model.SeatRepository {
	return &SeatRepository{
		dbClient: dbClient,
	}
}

func (r *SeatRepository) GetSeatsByMovieID(ctx context.Context, movieID int) ([]entities.Seat, error) {

	query := sq.Select("id", "movie_id", "seat_number", "status").
		From("seats").
		Where(sq.Eq{"movie_id": movieID}).
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

	var seats []entities.Seat
	for rows.Next() {
		var seat entities.Seat
		if err := rows.Scan(&seat.ID, &seat.MovieID, &seat.SeatNumber, &seat.Status); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seats, nil
}

func (r *SeatRepository) GetSeatByID(ctx context.Context, seatID int) (*entities.Seat, error) {

	query := sq.Select("id", "movie_id", "seat_number", "status").
		From("seats").
		Where(sq.Eq{"id": seatID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.dbClient.Pool.QueryRow(ctx, sql, args...)

	var seat entities.Seat
	if err := row.Scan(&seat.ID, &seat.MovieID, &seat.SeatNumber, &seat.Status); err != nil {
		return nil, err
	}

	return &seat, nil
}

func (r *SeatRepository) UpdateSeatStatus(ctx context.Context, tx pgx.Tx, seatID int, status string) error {

	query := sq.Update("seats").
		Set("status", status).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": seatID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	return err
}
