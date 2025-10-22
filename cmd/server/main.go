package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/a-h/templ"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/starfederation/datastar-go/datastar"
	"github.com/winkler1/dgoat/web/views"
)

var counter int

func handleErr(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Printf("Handler error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handleErr(func(w http.ResponseWriter, r *http.Request) error {
		return views.Home().Render(r.Context(), w)
	}))

	r.Get("/counter", handleErr(func(w http.ResponseWriter, r *http.Request) error {
		sse := datastar.NewSSE(w, r)
		return sse.PatchElementTempl(views.Counter(counter))
	}))

	r.Get("/increment", handleErr(func(w http.ResponseWriter, r *http.Request) error {
		counter++
		sse := datastar.NewSSE(w, r)
		return sse.PatchElementTempl(views.Counter(counter))
	}))

	hostPort := flag.String("hostport", "localhost:8080", "Host and port to run the server on")
	log.Printf("Starting server on %s", *hostPort)
	log.Fatal(http.ListenAndServe(*hostPort, r))
}
