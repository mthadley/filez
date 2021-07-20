package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mthadley/filez/internal/server"
)

func main() {
	var basePath string

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		basePath = args[0]
	} else {
		basePath = "."
	}

	baseDir, err := filepath.Abs(basePath)
	if err != nil {
		log.Fatal("Not a valid base directory", err)
	}

	base := os.DirFS(baseDir)
	server := server.NewServer(base)

	fmt.Printf("Serving folder %s at localhost:8080...\n\n", baseDir)

	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
