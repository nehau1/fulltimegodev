package db

const (
	DBNAME     = "hotel-reservation"
	DBURI      = "mongodb://localhost:27017"
	TestDBNAME = "hotel-reservation-test"
)

type Map map[string]any

type Pagination struct {
	Page  int64
	Limit int64
}

type Store struct {
	Hotel   HotelStore
	User    UserStore
	Room    RoomStore
	Booking BookingStore
}
