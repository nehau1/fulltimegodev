package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/Stiffjobs/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, from, till time.Time, ) *types.Booking {
	booking := &types.Booking{
		UserID: userID,
		RoomID: roomID,
		FromDate: from,
		TillDate: till,
		Canceled: false,
	}
	insertedBooking, err := store.Booking.Create(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func AddRoom(store *db.Store, size string, price float64, seaside bool, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size: size,
		Seaside: seaside,
		Price: price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.Create(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddHotel(store *db.Store,name, location string, rating int,  rooms []primitive.ObjectID) *types.Hotel {
	var roomIDs = rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Rating: rating,
		Rooms: roomIDs,
	}
	insertedHotel, err := store.Hotel.Create(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s%s@gmail.com", fn,ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.Create(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}
