package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// NewHandler initializes a new serve router
func NewHandler(s service) *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/{plugin}/{format}", getFeedHandler(s))
	})

	return r
}

// getUserHandler returns information about a given existing user
func getFeedHandler(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hit route")
		plugin := chi.URLParam(r, "plugin")
		format := chi.URLParam(r, "format")

		if plugin == "" {
			http.Error(w, errors.New("plugin not allowed to be empty").Error(), http.StatusInternalServerError)
			return
		}

		if format == "" {
			format = "rss"
		}
		s, err := s.Storage.Get(fmt.Sprintf("%s_%s", format, plugin))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(s))
	}
}
