package entrypoint

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thatpix3l/collge_event_website/src/api"
)

var db_url string = "postgres://postgres:postgres@127.0.0.1/college_event_website"

func Main() {

	log.SetFlags(log.Lshortfile)

	// Get database connection
	if conn, err := pgxpool.New(context.Background(), db_url); err != nil {
		log.Fatal(err)
	} else {
		// Store pool connection
		api.GlobalState.Pool = conn
	}
	defer api.GlobalState.Pool.Close()

	// Create new router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(api.Authentication)
	r.Use(api.Cleanup)

	// Add routes
	for path, pathFuncs := range api.HandleFuncs {
		for method, fn := range pathFuncs {
			r.MethodFunc(method, path, fn)
		}
	}

	http.ListenAndServe(":3000", r)
}
