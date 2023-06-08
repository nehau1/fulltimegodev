package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Stiffjobs/hotel-reservation/api/middleware"
	"github.com/Stiffjobs/hotel-reservation/db/fixtures"
	"github.com/Stiffjobs/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user           = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(tdb.Store, "hotel", "a", 4, nil)
		room           = fixtures.AddRoom(tdb.Store, "small", 4.4, true, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		jwtApp         = app.Group("/", middleware.JWTAuthentication(tdb.User))
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	jwtApp.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	token := CreateTokenFromUser(user)
	req.Header.Add("X-Api-Token", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	if bookingResp.ID != booking.ID {
		t.Fatalf("expected booking id to be %s, got %s", booking.ID.Hex(), bookingResp.ID.Hex())
	}

	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected booking user id to be %s, got %s", booking.UserID.Hex(), bookingResp.UserID.Hex())
	}
}

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser      = fixtures.AddUser(tdb.Store, "admin", "user", true)
		user           = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(tdb.Store, "hotel", "a", 4, nil)
		room           = fixtures.AddRoom(tdb.Store, "small", 4.4, true, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(tdb.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	_ = booking
	admin.Get("/", bookingHandler.HandleGetListBooking)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code to be %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s but got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s but got %s", booking.UserID, have.UserID)
	}

	// test normal user
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected status code to be %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
