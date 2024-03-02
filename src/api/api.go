package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/thatpix3l/collge_event_website/src/gen_sql"
	"github.com/thatpix3l/collge_event_website/src/utils"
)

// Parse coordinate from HTTP request; store and return copy of it.
func createCoordinate(queries *gen_sql.Queries, req *http.Request) (gen_sql.Coordinate, error) {

	// Empty inserted coordinate
	var inserted_coord gen_sql.Coordinate

	// Deserialize new coordinates
	var coordParams gen_sql.CreateCoordinateParams
	if err := json.NewDecoder(req.Body).Decode(&coordParams); err != nil {
		log.Println("Unable to deserialize coordinate values")
		return inserted_coord, err
	}

	// Store new record of coordinates into DB, get copy of what was inserted
	if coord, err := queries.CreateCoordinate(req.Context(), coordParams); err != nil {
		log.Println("Unable to create coordinate for university")
		return inserted_coord, err
	} else {
		inserted_coord = coord
	}

	return inserted_coord, nil
}

// Closure handler to create a university
func ClosureCreateUniversity(sharedState utils.SharedState) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {

		// Acquire database connection from pool
		conn, err := sharedState.Pool.Acquire(context.Background())
		if err != nil {
			log.Println("Unable to acquire database connection from pool")
			return
		}
		defer conn.Release()

		// Create queries connection
		queries := gen_sql.New(conn)

		// Create new coordinate
		inserted_coord, err := createCoordinate(queries, req)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Finished inserting coordinate: %s, %f, %f\n", inserted_coord.Title, inserted_coord.Latitude, inserted_coord.Longitude)

	}
}
