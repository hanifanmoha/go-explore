package movie_service

import (
	"context"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
	repositories_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/repositories"
	services_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/services"
)

type MovieService struct {
	movieRepo repositories_model.MovieRepository
}

func NewMovieService(movieRepo repositories_model.MovieRepository) services_model.MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

func (s *MovieService) GetAllMovies(ctx context.Context) ([]entities.Movie, error) {
	return s.movieRepo.GetAllMovies(ctx)
}

func (s *MovieService) GetMovieByID(ctx context.Context, id int) (*entities.Movie, error) {
	return s.movieRepo.GetMovieByID(ctx, id)
}
