package repositories_model

import (
	"context"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
	"github.com/jackc/pgx/v5"
)

type SeatRepository interface {
	GetSeatsByMovieID(ctx context.Context, movieID int) ([]entities.Seat, error)
	GetSeatByID(ctx context.Context, seatID int) (*entities.Seat, error)
	UpdateSeatStatus(ctx context.Context, tx pgx.Tx, seatID int, status string) error
}
