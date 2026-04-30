package movie_repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/clients"
	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
	repositories_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/repositories"
)

type MovieRepository struct {
	dbClient *clients.DatabaseClient
}

func NewMovieRepository(dbClient *clients.DatabaseClient) repositories_model.MovieRepository {
	return &MovieRepository{
		dbClient: dbClient,
	}
}

func (r *MovieRepository) GetAllMovies(ctx context.Context) ([]entities.Movie, error) {

	query := sq.Select("id", "title", "description", "image_url").
		From("movies").
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

	var movies []entities.Movie
	for rows.Next() {
		var movie entities.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ImageURL); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) GetMovieByID(ctx context.Context, id int) (*entities.Movie, error) {

	query := sq.Select("id", "title", "description", "image_url").
		From("movies").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.dbClient.Pool.QueryRow(ctx, sql, args...)

	var movie entities.Movie
	if err := row.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ImageURL); err != nil {
		return nil, err
	}

	return &movie, nil
}
