package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/Stiffjobs/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func postTestUser(t *testing.T, userStore db.UserStore) *types.User{
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "helloworld@gmail.com",
		Password:  "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	user, err = userStore.Create(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	testdb := setup(t)
	defer testdb.teardown(t)
	insertedUser:= postTestUser(t, testdb)

	app := fiber.New()
	authHandler := NewAuthHandler(testdb)
	app.Post("/auth", authHandler.HandleAuthenticate)


	params := AuthParams{
		Email: "helloworld@gmail.com",
		Password: "password",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	t.Logf("resp status: %v", resp.StatusCode)
	if err != nil {
		t.Fatal(err)
	}
	var authResponse AuthResponse 
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatal(err)
	} 
	t.Log(authResponse)
	if authResponse.Token == "" {
		t.Fatal("expected token to be present")
	}

	//Set the encrypted password to empty string.
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResponse.User) {
		fmt.Println(insertedUser)
		fmt.Println(authResponse.User)
		t.Fatal("expected user to be equal")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {

	testdb := setup(t)
	defer testdb.teardown(t)
	_ = postTestUser(t, testdb)

	app := fiber.New()
	authHandler := NewAuthHandler(testdb)
	app.Post("/auth", authHandler.HandleAuthenticate)
	params := AuthParams{
		Email: "helloworld@gmail.com",
		Password: "wrongpassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest{
		t.Fatal("expected status code to be 400")
	}
	var genResp genericResp

	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
	}

	if genResp.Message != "invalid credentials" {
		t.Fatalf("expected gen response msg to be <invalid credentials> but got %s", genResp.Message)
	}

}