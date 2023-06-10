package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/storage/bbolt"
)

func main() {
	
	store := bbolt.New(bbolt.Config{
		Database: "times.db",
		Bucket:   "times",
		Reset:    false,
	})
	defer store.Close()

	times := map[string]string{
		"ate": "0",
		"ate_previous": "0",
		"pee": "0",
		"pee_previous": "0",
		"poop": "0",
		"poop_previous": "0",
	}
	exp := 0*time.Second

	for key, val := range times {
		err := store.Set(key, []byte(val), exp)
		if err != nil {
			log.Println("Error setting value:", err)
		}
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/time/:timerType", func(c *fiber.Ctx) error {
		timerType := c.Params("timerType")
		byteValue, _ := store.Get(c.Params("timerType"))
		value := string(byteValue)
		log.Println(timerType,"is",value)
		ret := struct{
			Time string
		}{Time: value}
		return c.JSON(ret)
	})

	app.Post("/time/:timerType/:timeStamp", func(c *fiber.Ctx) error {
		timerType := c.Params("timerType")
		timeStamp := c.Params("timeStamp")
		log.Println("update:", timerType, "with value:", timeStamp)
		store.Set(timerType, []byte(timeStamp), exp)
		return c.JSON(struct{Time string}{timeStamp})
	})

	log.Println("live")
	log.Fatal(app.Listen(":3000"))
}
