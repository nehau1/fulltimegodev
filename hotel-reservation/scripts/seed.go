package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Stiffjobs/hotel-reservation/api"
	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/Stiffjobs/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func main() {
	ctx := context.Background()
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User: db.NewMongoUserStore(client),
		Hotel: hotelStore,
		Room: db.NewMongoRoomStore(client,hotelStore),
		Booking: db.NewMongoBookingStore(client),
	}
	john := fixtures.AddUser(store, "john", "doe", false)

	fmt.Printf("john -> %s\n", api.CreateTokenFromUser(john))
	adminUser := fixtures.AddUser(store,"admin", "user", true)
	fmt.Printf("admin -> %s\n", api.CreateTokenFromUser(adminUser))
	fixtures.AddHotel(store, "Bellagio", "Las Vegas", 5,nil)
	fixtures.AddHotel(store, "the cozy hotel", "The Netherlands", 4,nil)
	hotel := fixtures.AddHotel(store, "Bellucia", "France", 3, nil)
	smallRoom := fixtures.AddRoom(store, "small", 89.99, true, hotel.ID)
	mediumRoom := fixtures.AddRoom(store, "medium", 289.99, true, hotel.ID)
	largeRoom := fixtures.AddRoom(store, "large",  389.99, false, hotel.ID)
	fixtures.AddBooking(store,adminUser.ID, smallRoom.ID, time.Now().AddDate(0, 1, 0), time.Now().AddDate(0, 1, 1))
	fixtures.AddBooking(store, adminUser.ID, mediumRoom.ID, time.Now().AddDate(0, 1, 15), time.Now().AddDate(0, 2, 1))
	booking := fixtures.AddBooking(store, adminUser.ID, largeRoom.ID, time.Now().AddDate(0, 2, 20), time.Now().AddDate(0, 3, 1))
	fmt.Println(booking)
}

