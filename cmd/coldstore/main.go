package main

import (
	"log"
	"net/http"
	"os"

	"coldstore/internal/console"
	"coldstore/internal/service"
)

func main() {
	dataPath := os.Getenv("COLDSTORE_DATA")
	if dataPath == "" {
		dataPath = "coldstore-data.json"
	}
	srv, err := service.New(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.Recover(); err != nil {
		log.Fatal(err)
	}
	addr := os.Getenv("COLDSTORE_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("coldstore listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, console.NewRouter(srv)))
}
