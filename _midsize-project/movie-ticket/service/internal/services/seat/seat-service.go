package seat_service

import (
	"context"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/entities"
	repositories_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/repositories"
	services_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/services"
)

type SeatService struct {
	seatRepo repositories_model.SeatRepository
}

func NewSeatService(seatRepo repositories_model.SeatRepository) services_model.SeatService {
	return &SeatService{
		seatRepo: seatRepo,
	}
}

func (s *SeatService) GetSeatsByMovieID(ctx context.Context, movieID int) ([]entities.Seat, error) {
	return s.seatRepo.GetSeatsByMovieID(ctx, movieID)
}
