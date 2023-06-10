package main

import (
	"log"
	"time"
	
    "embed"
    "net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/storage/bbolt"
)

// Embed a directory
//go:embed dist/*
var embedDirStatic embed.FS


func main() {
	
	store := bbolt.New(bbolt.Config{
		Database: "times.db",
		Bucket:   "times",
		Reset:    false,
	})
	defer store.Close()

	exp := 0*time.Second
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

	app.Use("/", filesystem.New(filesystem.Config{
        Root: http.FS(embedDirStatic),
        PathPrefix: "dist",
        Browse: true,
    }))

	log.Println("live")
	log.Fatal(app.Listen(":3000"))
}
