package main

import (
	"net"
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

	cache := cache.New(time.Duration(config.CACHE_EXPIRATION)*time.Minute, time.Duration(config.CACHE_EXPIRED_PURGE)*time.Minute)
	storageRepo, err := store.NewMemRepository(cache)
	if err != nil {
		return
	}

	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	c := &http.Client{
		Timeout:   time.Second * 10,
		Transport: t,
	}

	pluginRepo := repository.NewMemRepo()
	pluginRepo.Install(scmp.NewPlugin(l, c))

	runner := runner.NewRunner(l, pluginRepo, storageRepo, config.REFRESH_INTERVAL)
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
