package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mthadley/filez/internal/server"
)

func main() {
	baseDir := getBaseDir()
	base := os.DirFS(baseDir)

	port := getPort()

	fmt.Printf("Serving folder %s at localhost:%d...\n\n", baseDir, port)

	server := server.NewServer(base)

	err := http.ListenAndServe(":"+strconv.Itoa(port), server)
	if err != nil {
		log.Fatal(err)
	}
}

func getPort() int {
	port := 8080

	rawEnvPort, envPortPresent := os.LookupEnv("PORT")
	if envPortPresent {
		envPort, err := strconv.Atoi(rawEnvPort)
		if err == nil {
			port = envPort
		}
	}

	return *flag.Int("p", port, "The port to listen on.")
}

func getBaseDir() string {
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

	return baseDir
}
