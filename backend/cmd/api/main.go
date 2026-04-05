package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/jonnie-z/notes-app/internal/httpapi"
	"github.com/jonnie-z/notes-app/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: '%v'", err)
	}

	args := os.Args[1:]
	var storeType store.StoreType

	if len(args) == 0 {
		storeType = store.StoreJSON
	} else {
		switch args[0] {
		case "json":
			storeType = store.StoreJSON
		case "mem":
			storeType = store.StoreInMemory
		case "sql":
			storeType = store.StoreSQL
		}
	}

	app := newApp(storeType)
	api := &httpapi.API{App: app}

	mux := api.Routes()

	fmt.Println("Starting Server on", app.Port)
	if err := http.ListenAndServe(app.Port, mux); err != nil {
		log.Fatal(err)
	}
}

