package main

import entrypoint "github.com/thatpix3l/collge_event_website/src"

//go:generate jet -dsn postgresql://postgres:postgres@127.0.0.1/college_event_website?sslmode=disable -schema cew -path src/gen_sql
//go:generate templ generate -path src/gen_templ

func main() {
	entrypoint.Main()
}
