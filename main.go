package main

import entrypoint "github.com/thatpix3l/collge_event_website/src"

//go:generate sqlc generate
//go:generate templ generate -path src/gen_templ
//go:generate easyjson -all src/gen_sql/queries.sql.go

func main() {
	entrypoint.Main()
}
