package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xomatix/silly-syntax-backend-bonanza/api"
	"github.com/xomatix/silly-syntax-backend-bonanza/database"
	pluginmanager "github.com/xomatix/silly-syntax-backend-bonanza/pluginManager"
	_ "github.com/xomatix/silly-syntax-backend-bonanza/statik"

	"github.com/rakyll/statik/fs"
)

func main() {

	mux := http.NewServeMux()

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.FileServer(statikFS))

	mux.HandleFunc("/bonanza", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🍌🍌🍌"))
	})

	database.InitDatabase()
	database.LoadTablesConfig()
	database.InitDatabasePermissions()
	database.InitDatabaseViews()
	api.InitApiRoutes(mux)

	pluginmanager.GetPluginLoader().LoadPlugins(mux)

	handler := enableCORS(mux)

	fmt.Println("Starting server on http://localhost:8080")
	fmt.Println(http.ListenAndServe(":8080", handler))

}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
