package classbooking

type ClassBookingID string

type BookingStatus string

const (
	BookingStatusBooked    BookingStatus = "BOOKED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusAttended  BookingStatus = "ATTENDED"
)

func (id ClassBookingID) String() string {
	return string(id)
}

func validateBookingStatus(status BookingStatus) error {
	switch status {
	case BookingStatusBooked, BookingStatusCancelled, BookingStatusAttended:
		return nil
	default:
		return ErrInvalidStatus
	}
}
