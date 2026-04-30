package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	services_model "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/models/services"
)

type Handler struct {
	bookingService services_model.BookingService
	movieService   services_model.MovieService
	seatService    services_model.SeatService
}

func NewHandler(bookingService services_model.BookingService, movieService services_model.MovieService, seatService services_model.SeatService) *Handler {
	return &Handler{
		bookingService: bookingService,
		movieService:   movieService,
		seatService:    seatService,
	}
}

type MovieResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type SeatResponse struct {
	ID         int    `json:"id"`
	MovieID    int    `json:"movie_id"`
	SeatNumber string `json:"seat_number"`
	Status     string `json:"status"`
}

func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {

	movies, err := h.movieService.GetAllMovies(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}

	var response []MovieResponse
	for _, movie := range movies {
		response = append(response, MovieResponse{
			ID:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			ImageURL:    movie.ImageURL,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {

	_id := r.PathValue("id")
	id, err := strconv.Atoi(_id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie, err := h.movieService.GetMovieByID(r.Context(), id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to fetch movie", http.StatusInternalServerError)
		return
	}

	response := MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		ImageURL:    movie.ImageURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) GetSeats(w http.ResponseWriter, r *http.Request) {

	_movieId := r.PathValue("id")
	movieId, err := strconv.Atoi(_movieId)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	seats, err := h.seatService.GetSeatsByMovieID(r.Context(), movieId)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to fetch seats", http.StatusInternalServerError)
		return
	}

	var response []SeatResponse
	for _, seat := range seats {
		response = append(response, SeatResponse{
			ID:         seat.ID,
			MovieID:    seat.MovieID,
			SeatNumber: seat.SeatNumber,
			Status:     seat.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
