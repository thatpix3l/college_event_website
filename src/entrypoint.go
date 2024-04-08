package entrypoint

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thatpix3l/collge_event_website/src/api"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db_url string = "postgres://postgres:postgres@127.0.0.1/college_event_website"

func Main() {

	log.SetFlags(log.Lshortfile)

	// Get DB connection
	db, err := sql.Open("pgx", db_url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api.GlobalState.Db = db

	// Create new router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(api.Authentication)

	// Add routes
	for path, pathFuncs := range api.HandleFuncs {
		for method, fn := range pathFuncs {
			r.MethodFunc(method, path, fn)
		}
	}

	http.ListenAndServe(":3000", r)
}
