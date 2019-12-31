package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fwojciec/gqlgen-sqlc-example/gqlgen" // update the username
	"github.com/fwojciec/gqlgen-sqlc-example/pg"     // update the username
)

func main() {
	// initialize the db
	db, err := pg.Open("dbname=gqlgen_sqlc_example_db sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// initialize the repository
	repo := pg.NewRepository(db)

	// configure the server
	mux := http.NewServeMux()
	mux.Handle("/", gqlgen.NewPlaygroundHandler("/query"))
	mux.Handle("/query", gqlgen.NewHandler(repo))

	// run the server
	port := ":8080"
	fmt.Fprintf(os.Stdout, "🚀 Server ready at http://localhost%s\n", port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(port, mux))
}
