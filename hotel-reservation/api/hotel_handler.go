package api

import (
	"errors"

	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "hotel not found"})
		}
		return err
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetListHotel(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetList(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetListRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelID": oid}

	rooms, err := h.store.Room.GetList(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
