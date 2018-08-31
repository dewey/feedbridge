package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// NewHandler initializes a new feed router
func NewHandler(s service) *chi.Mux {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/list", getPluginListHandler(s))
		r.Get("/{plugin}/{format}", getFeedHandler(s))
	})

	r.Group(func(r chi.Router) {
		r.Use(s.authenticator)
		r.Post("/{plugin}/refresh", refreshFeedHandler(s))
	})

	return r
}

// authenticator checks if the users sends the correct token to access the internal routes
func (s *service) authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("auth_token") == "" || s.cfg.APIToken == "" || q.Get("auth_token") != s.cfg.APIToken {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// getFeedHandler returns a feed in a specified format
func getFeedHandler(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plugin := chi.URLParam(r, "plugin")
		if plugin == "" {
			http.Error(w, errors.New("plugin not allowed to be empty").Error(), http.StatusInternalServerError)
			return
		}

		format := chi.URLParam(r, "format")
		if format == "" {
			format = "rss"
		}
		s, err := s.ServeFeed(format, plugin)
		if err != nil {
			http.Error(w, errors.New("there was an error serving the feed").Error(), http.StatusInternalServerError)
			return
		}

		switch format {
		case "atom":
			w.Header().Set("Content-Type", "application/atom+xml")
		case "json":
			w.Header().Set("Content-Type", "application/json")
		default:
			w.Header().Set("Content-Type", "application/rss+xml")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(s))
	}
}

// getPluginListHandler returns a list of all available plugins
func getPluginListHandler(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plugins := s.ListFeeds()
		b, err := json.Marshal(plugins)
		if err != nil {
			http.Error(w, errors.New("there was an error listing the plugins").Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

// refreshFeedHandler returns a list of all available plugins
func refreshFeedHandler(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		plugin := chi.URLParam(r, "plugin")
		if plugin == "" {
			http.Error(w, errors.New("plugin not allowed to be empty").Error(), http.StatusInternalServerError)
			return
		}
		if err := s.RefreshFeed(plugin); err != nil {
			http.Error(w, errors.New("there was an error listing the plugins").Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
