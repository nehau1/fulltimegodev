package api

import (
	"context"
	"testing"

	"github.com/Stiffjobs/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdburi = "mongodb://localhost:27017"
	dbname = "hotel-reservation-test"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown (t *testing.T){
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func  setup (t *testing.T) *testdb{
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		t.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		Store: &db.Store{
			User: db.NewMongoUserStore(client),
			Hotel: hotelStore,
			Room: db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}