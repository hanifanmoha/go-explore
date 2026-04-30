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

type BookingResponse struct {
	ID      int            `json:"id"`
	MovieID int            `json:"movie_id"`
	SeatID  int            `json:"seat_id"`
	UserID  string         `json:"user_id"`
	Status  string         `json:"status"`
	Movie   *MovieResponse `json:"movie,omitempty"`
	Seat    *SeatResponse  `json:"seat,omitempty"`
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

func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {

	userID := r.PathValue("user_id")
	if userID == "" {
		slog.Error("user_id query parameter is required")
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	bookings, err := h.bookingService.GetBookingsByUserID(r.Context(), userID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	var response []BookingResponse
	for _, booking := range bookings {
		response = append(response, BookingResponse{
			ID:      booking.ID,
			MovieID: booking.MovieID,
			SeatID:  booking.SeatID,
			UserID:  booking.UserID,
			Status:  booking.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {

	type CreateBookingRequest struct {
		UserID  string `json:"user_id"`
		MovieID int    `json:"movie_id"`
		SeatID  int    `json:"seat_id"`
	}

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.bookingService.OrderSeat(r.Context(), req.MovieID, req.SeatID, req.UserID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking created successfully"})

}
