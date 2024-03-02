package entrypoint

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thatpix3l/collge_event_website/src/api"
	app "github.com/thatpix3l/collge_event_website/src/gen_templ"
	"github.com/thatpix3l/collge_event_website/src/utils"
)

var db_url string = "postgres://postgres:postgres@127.0.0.1/college_event_website"

func apiPath(subpath string) string {
	return "/api/v1/" + subpath
}

func Main() {

	sharedState := utils.SharedState{}

	// Get database connection
	if conn, err := pgxpool.New(context.Background(), db_url); err != nil {
		log.Fatal(err)
	} else {
		sharedState.Pool = conn
	}
	defer sharedState.Pool.Close()

	// Create new router; add handlers to it.
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post(apiPath("create_university"), api.ClosureCreateUniversity(sharedState))
	r.Handle("/", templ.Handler(app.Main()))

	http.ListenAndServe(":3000", r)
}
