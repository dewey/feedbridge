package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/caarlos0/env"
	"github.com/gobuffalo/packr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dewey/feedbridge/api"
	"github.com/dewey/feedbridge/plugin"
	"github.com/dewey/feedbridge/plugins/scmp"
	"github.com/dewey/feedbridge/runner"
	"github.com/dewey/feedbridge/store"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	cache "github.com/patrickmn/go-cache"
)

func main() {
	var config struct {
		RefreshInterval   int    `env:"REFRESH_INTERVAL" envDefault:"15"`
		CacheExpiration   int    `env:"CACHE_EXPIRATION" envDefault:"30"`
		CacheExpiredPurge int    `env:"CACHE_EXPIRED_PURGE" envDefault:"60"`
		Environment       string `env:"ENVIRONMENT" envDefault:"develop"`
		Port              int    `env:"PORT" envDefault:"8080"`
		APIHostname       string `env:"API_HOSTNAME" envDefault:"http://localhost"`
	}
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	l := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	switch strings.ToLower(config.Environment) {
	case "develop":
		l = level.NewFilter(l, level.AllowInfo())
	case "prod":
		l = level.NewFilter(l, level.AllowError())
	}
	l = log.With(l, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cache := cache.New(time.Duration(config.CacheExpiration)*time.Minute, time.Duration(config.CacheExpiredPurge)*time.Minute)
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

	pluginRepo := plugin.NewMemRepo()
	pluginRepo.Install(scmp.NewPlugin(l, c))

	runner := runner.NewRunner(l, pluginRepo, storageRepo, config.RefreshInterval)
	go runner.Start()

	apiService := api.NewService(l, storageRepo, pluginRepo)

	templates := packr.NewBox("./ui/templates")
	assets := packr.NewBox("./ui/assets")

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index.tmpl").Parse(templates.String("index.tmpl"))
		if err != nil {
			http.Error(w, errors.New("couldn't serve template").Error(), http.StatusInternalServerError)
			return
		}
		if err := t.Execute(w, apiService.ListFeeds()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(assets)))
	r.Handle("/metrics", promhttp.Handler())

	// TODO(dewey): Switch to promhttp middleware instead of this deprecated one
	r.Mount("/feed", prometheus.InstrumentHandler("feed", api.NewHandler(*apiService)))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("nothing to see here"))
	})

	l.Log("msg", fmt.Sprintf("feedbridge listening on %s:%d", config.APIHostname, config.Port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
	if err != nil {
		panic(err)
	}
}
