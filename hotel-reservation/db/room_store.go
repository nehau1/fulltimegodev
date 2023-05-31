package db

import (
	"context"

	"github.com/Stiffjobs/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	Create(context.Context, *types.Room) (*types.Room, error)
	GetList(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(DBNAME).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) Create(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)

	// update the hotel with this room id
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}

	s.HotelStore.Update(ctx, filter, update)

	return room, nil
}

func (s *MongoRoomStore) GetList(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
