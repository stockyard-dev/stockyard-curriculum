package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stockyard-dev/stockyard-curriculum/internal/server"
	"github.com/stockyard-dev/stockyard-curriculum/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9813"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./curriculum-data"
	}

	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("curriculum: %v", err)
	}
	defer db.Close()

	srv := server.New(db, server.DefaultLimits())

	fmt.Printf("\n  Curriculum — Self-hosted lesson planning and curriculum tracking\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Questions? hello@stockyard.dev — I read every message\n\n", port, port)
	log.Printf("curriculum: listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
