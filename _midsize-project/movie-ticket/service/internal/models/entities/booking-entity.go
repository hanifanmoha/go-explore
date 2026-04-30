package entities

const (
	BOOKING_STATUS_PENDING   = "pending"
	BOOKING_STATUS_CONFIRMED = "confirmed"
	BOOKING_STATUS_CANCELLED = "cancelled"
)

type Booking struct {
	ID      int
	SeatID  int
	UserID  string
	MovieID int
	Status  string
}
