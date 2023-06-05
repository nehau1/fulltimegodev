package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Stiffjobs/hotel-reservation/api"
	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/Stiffjobs/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := userStore.Create(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedRoom(size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := roomStore.Create(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}

	insertedBooking, err := bookingStore.Create(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Booking ID: %s\n", insertedBooking.ID.String())
	return insertedBooking
}

func seedHotel(name, location string, rating int) *types.Hotel {

	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	insertedHotel, err := hotelStore.Create(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel

}

func main() {
	seedUser(false, "John", "Doe", "helloworld@gmail.com", "password")
	adminUser := seedUser(true, "admin", "admin", "admin@gmail.com", "adminpassword")
	seedHotel("Bellagio", "Las Vegas", 5)
	seedHotel("the cozy hotel", "The Netherlands", 4)
	hotel := seedHotel("Bellucia", "France", 3)
	smallRoom := seedRoom("small", true, 89.99, hotel.ID)
	mediumRoom := seedRoom("medium", true, 289.99, hotel.ID)
	largeRoom := seedRoom("large", false, 389.99, hotel.ID)
	seedBooking(adminUser.ID, smallRoom.ID, time.Now().AddDate(0, 1, 0), time.Now().AddDate(0, 1, 1))
	seedBooking(adminUser.ID, mediumRoom.ID, time.Now().AddDate(0, 1, 15), time.Now().AddDate(0, 2, 1))
	seedBooking(adminUser.ID, largeRoom.ID, time.Now().AddDate(0, 2, 20), time.Now().AddDate(0, 3, 1))
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	userStore = db.NewMongoUserStore(client)
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	bookingStore = db.NewMongoBookingStore(client)
}
