package main

import (
	"log"
	"net/http"
	"os"

	"grps-go-redis-psql/internal/app"
)

func main() {
	application := app.New()
	defer application.Close()
	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}

	log.Println("server started on port", port)

	if err := http.ListenAndServe(":"+port, application.Router()); err != nil {
		log.Fatal(err)
	}
}
