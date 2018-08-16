package api

import (
	"encoding/json"
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
		r.Get("/list", getPluginListHandler(s))
		r.Get("/{plugin}/{format}", getFeedHandler(s))
	})

	return r
}

// getUserHandler returns information about a given existing user
func getFeedHandler(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plugin := chi.URLParam(r, "plugin")
		format := chi.URLParam(r, "format")

		if plugin == "" {
			http.Error(w, errors.New("plugin not allowed to be empty").Error(), http.StatusInternalServerError)
			return
		}

		if format == "" {
			format = "rss"
		}
		s, err := s.StorageRepository.Get(fmt.Sprintf("%s_%s", format, plugin))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(s))
	}
}

// getPluginListHandler returns a list of all available plugins
func getPluginListHandler(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := s.PluginRepository.All()
		var plugins []string
		for _, p := range p {
			plugins = append(plugins, p.String())
		}
		b, err := json.Marshal(plugins)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
