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

	mux.Handle("/", spaHandler("./build", "index.html"))

	fmt.Println("Starting Server on", app.Port)
	if err := http.ListenAndServe(app.Port, mux); err != nil {
		log.Fatal(err)
	}
}

func spaHandler(staticPath, indexPath string) http.Handler {
    fs := http.FileServer(http.Dir(staticPath))

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        path := staticPath + r.URL.Path

        _, err := os.Stat(path)
        if os.IsNotExist(err) {
            http.ServeFile(w, r, staticPath+"/"+indexPath)
            return
        }

        fs.ServeHTTP(w, r)
    })
}