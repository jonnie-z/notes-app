package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jonnie-z/notes-app/internal/app"
	"github.com/jonnie-z/notes-app/internal/httpapi"
	"github.com/jonnie-z/notes-app/internal/store"
)

const PORT = ":8080"

// const DATA_FILE = "./db.json"
// const TEMP_DATA_FILE = DATA_FILE + ".tmp"



func main() {
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
		}
	}

	app := app.NewApp(storeType)
	api := &httpapi.API{App: app}

	mux := api.Routes()

	fmt.Println("Starting Server on", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}
}

