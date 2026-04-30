package services_model

import (
	"context"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
)

type MovieService interface {
	GetAllMovies(ctx context.Context) ([]entities.Movie, error)
	GetMovieByID(ctx context.Context, id int) (*entities.Movie, error)
}
