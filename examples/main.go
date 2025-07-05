package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gocss-example/templates"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Serve static files (CSS)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Render Templ component
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		component := templates.Hello("GOCSS Example")
		component.Render(r.Context(), w)
	})

	fmt.Println("Server starting on :3000")
	http.ListenAndServe(":3000", r)
}
