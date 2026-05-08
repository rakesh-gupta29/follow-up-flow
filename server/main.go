package main

import (
	"log"

	"github.com/shingo/server/bootstrap"
)

func main() {
	app := bootstrap.New()
	defer app.DB.Close()

	log.Fatal(app.Server.Listen(":" + app.Config.AppPort))
}
