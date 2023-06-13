package api

import (
	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
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
		return ErrInvalidID()
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetListHotel(c *fiber.Ctx) error {
	var query HotelQueryParams
	if err := c.QueryParser(&query); err != nil {
		return ErrBadRequest()
	}
	filter := db.Map{
		"rating": query.Rating,
	}
	hotels, err := h.store.Hotel.GetList(c.Context(), filter, &query.Pagination)
	if err != nil {
		return err
	}
	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(query.Page),
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetListRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := db.Map{"hotelID": oid}

	rooms, err := h.store.Room.GetList(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
