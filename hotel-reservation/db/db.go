package db

const (
	DBNAME     = "hotel-reservation"
	DBURI      = "mongodb://localhost:27017"
	TestDBNAME = "hotel-reservation-test"
)

type Store struct {
	Hotel HotelStore
	User  UserStore
	Room  RoomStore
}
