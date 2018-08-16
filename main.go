package main

import (
	"net/http"
	"os"
	"time"

	"github.com/dewey/feedbridge/api"
	"github.com/dewey/feedbridge/plugins/scmp"
	"github.com/dewey/feedbridge/repository"
	"github.com/dewey/feedbridge/runner"
	"github.com/dewey/feedbridge/store"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/kelseyhightower/envconfig"
	cache "github.com/patrickmn/go-cache"
)

// Config contains all configuration options that can be overwritten with environment variables
type Config struct {
	REFRESH_INTERVAL    int `default:"10"`
	CACHE_EXPIRATION    int `default:"15"`
	CACHE_EXPIRED_PURGE int `default:"20"`
}

func main() {
	var config Config
	err := envconfig.Process("fb", &config)
	if err != nil {
		panic(err)
	}
	l := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	l = log.With(l, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	c := cache.New(time.Duration(config.CACHE_EXPIRATION)*time.Minute, time.Duration(config.CACHE_EXPIRED_PURGE)*time.Minute)
	storageRepo, err := store.NewMemRepository(c)
	if err != nil {
		return
	}

	pluginRepo := repository.NewMemRepo()
	pluginRepo.Install(scmp.NewPlugin(l))

	runner := runner.NewRunner(l, config.REFRESH_INTERVAL, pluginRepo, storageRepo)
	go runner.Start()

	apiService := api.NewService(storageRepo, pluginRepo)

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("feedbridge"))
	})
	r.Mount("/feed", api.NewHandler(*apiService))
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
