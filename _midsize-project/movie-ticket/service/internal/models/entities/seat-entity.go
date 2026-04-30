package entities

const (
	SEAT_STATUS_AVAILABLE = "available"
	SEAT_STATUS_RESERVED  = "booked"
	SEAT_STATUS_SOLD      = "sold"
)

type Seat struct {
	ID         int
	MovieID    int
	SeatNumber string
	Status     string
}
