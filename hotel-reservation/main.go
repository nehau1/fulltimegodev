package main

import (
	"context"
	"flag"
	"log"

	"github.com/Stiffjobs/hotel-reservation/api"
	"github.com/Stiffjobs/hotel-reservation/api/middleware"
	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{
			"error": err.Error(),
		})

	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server.")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}

	//initialization
	var (
		app         = fiber.New(config)
		apiv1       = app.Group("/api/v1", middleware.JWTAuthentication)
		apiv1Public = app.Group("/api")
		hotelStore  = db.NewMongoHotelStore(client)
		roomStore   = db.NewMongoRoomStore(client, hotelStore)
		userStore   = db.NewMongoUserStore(client)
		userHandler = api.NewUserHandler(userStore)
		authHandler = api.NewAuthHandler(userStore)
		store       = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		hotelHandler = api.NewHotelHandler(store)
	)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Hello World",
		})
	})

	//auth handlers
	apiv1Public.Post("/auth", authHandler.HandleAuthenticate)

	//Versioned
	//user handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetListHotel)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotelByID)
	apiv1.Get("/hotel/:id/room", hotelHandler.HandleGetListRoom)
	log.Fatal(app.Listen(*listenAddr))
}
