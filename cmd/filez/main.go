package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mthadley/filez/internal/server"
)

func main() {
	baseDir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal("Not a valid base directory", err)
	}

	base := os.DirFS(baseDir)
	server := server.NewServer(base)

	fmt.Println("Starting server...")

	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
