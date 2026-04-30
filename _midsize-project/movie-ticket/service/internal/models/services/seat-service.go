package services_model

import (
	"context"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
)

type SeatService interface {
	GetSeatsByMovieID(ctx context.Context, movieID int) ([]entities.Seat, error)
}
