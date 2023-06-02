package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/Stiffjobs/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) Validate() error {
	now := time.Now()

	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book room in the past")
	}

	if p.FromDate.After(p.TillDate) {
		return fmt.Errorf("fromDate must be before tillDate")
	}

	if p.NumPersons < 1 {
		return fmt.Errorf("numPersons must be greater than 0")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{store: store}
}

func (h *RoomHandler) HandleGetListRoom(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetList(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type:    "error",
			Message: err.Error(),
		})
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type:    "error",
			Message: "Internal server error",
		})
	}

	ok, err = h.isRoomAvailableForBooking(c, roomID, params)

	if err != nil {
		return err
	}

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type:    "error",
			Message: fmt.Sprintf("Room %s is not available for booking", c.Params("id")),
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}
	fmt.Println(booking)
	inserted, err := h.store.Booking.Create(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(c *fiber.Ctx, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}

	bookings, err := h.store.Booking.GetList(c.Context(), where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0

	return ok, nil
}
