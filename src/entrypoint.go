package entrypoint

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thatpix3l/collge_event_website/src/api"
	"github.com/thatpix3l/collge_event_website/src/utils"
)

var db_url string = "postgres://postgres:postgres@127.0.0.1/college_event_website"

func Main() {

	log.SetFlags(log.Lshortfile)

	// Create instance of handlers
	h := api.Handlers{}

	// Get database connection
	if conn, err := pgxpool.New(context.Background(), db_url); err != nil {
		log.Fatal(err)
	} else {
		// Store pool connection
		h.Pool = conn
	}
	defer h.Pool.Close()

	// Create new router; add handlers to it.
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	// r.Use(h.Authentication)

	// Routes
	h.Add(r.Get, "/", h.ReadHome)
	h.Add(r.Post, utils.ApiPath("university"), h.CreateUniversity)
	h.Add(r.Get, utils.ApiPath("university"), h.ReadUniversities)
	h.Add(r.Post, utils.ApiPath("login"), h.CreateLogin)
	h.Add(r.Get, utils.ApiPath("login"), h.ReadLogin)
	h.Add(r.Post, utils.ApiPath("signup"), h.CreateStudent)
	h.Add(r.Get, utils.ApiPath("signup"), h.ReadSignup)

	http.ListenAndServe(":3000", r)
}
