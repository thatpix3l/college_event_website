package entrypoint

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alexflint/go-arg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thatpix3l/cew/src/api"
	"github.com/thatpix3l/cew/src/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Main() {

	log.SetFlags(log.Lshortfile)

	// Parse CLI arguments
	arg.MustParse(&config.Root)

	// Get DB connection
	db, err := sql.Open("pgx", config.Root.Dsn())
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

	http.ListenAndServe(config.Root.Host(), r)
}
