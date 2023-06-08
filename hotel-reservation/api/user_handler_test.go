package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Stiffjobs/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	testdb := setup(t)
	defer testdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(testdb.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "helloworld@gmail.com",
		FirstName: "James",
		LastName:  "Foo Hello world",
		Password:  "aisdufoiajs",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	t.Log(resp.StatusCode)
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Error(err)
	}

	if len(user.ID) == 0 {
		t.Errorf("expected a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected the EncryptedPassword not to be included in the json response")
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastName %s but got %s", params.LastName, user.LastName)
	}
}
