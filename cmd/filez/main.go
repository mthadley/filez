package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/mthadley/filez/internal/server"
)

func main() {
	baseDir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal("Not a valid base directory", err)
	}

	server := server.NewServer(baseDir)

	fmt.Println("Starting server...")

	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
