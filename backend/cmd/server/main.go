package main

import (
	"log"
	"os"

	"github.com/texas-holdem/backend/internal/api"
)

func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	} else if addr[0] != ':' {
		addr = ":" + addr
	}
	srv := api.New()
	if err := srv.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
