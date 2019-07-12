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
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dewey/feedbridge/api"
	"github.com/dewey/feedbridge/config"
	"github.com/dewey/feedbridge/plugin"
	"github.com/dewey/feedbridge/plugins/racefansnet"
	"github.com/dewey/feedbridge/plugins/roadsandkingdoms"
	"github.com/dewey/feedbridge/plugins/scmp"
	"github.com/dewey/feedbridge/runner"

	"github.com/dewey/feedbridge/store"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	l := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	switch strings.ToLower(cfg.Environment) {
	case "develop":
		l = level.NewFilter(l, level.AllowInfo())
	case "prod":
		l = level.NewFilter(l, level.AllowError())
	}
	l = log.With(l, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	c := &http.Client{
		Timeout:   time.Second * 15,
		Transport: t,
	}

	pluginRepo := plugin.NewMemRepo()
	pluginRepo.Install(scmp.NewPlugin(l, c))
	pluginRepo.Install(roadsandkingdoms.NewPlugin(l, c))
	pluginRepo.Install(racefansnet.NewPlugin(l, c))

	storageRepo, err := store.NewStoreBackend(cfg)
	if err != nil {
		return
	}

	runner := runner.NewRunner(l, pluginRepo, storageRepo, cfg.RefreshInterval)
	// No scheduled scrapes running in development, use the refresh route to test plugins
	if cfg.Environment != "develop" {
		go runner.Start()
	}

	apiService := api.NewService(l, cfg, storageRepo, pluginRepo, runner)

	templates := packr.NewBox("../../ui/templates")
	assets := packr.NewBox("../../ui/assets")

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index.tmpl").Parse(templates.String("index.tmpl"))
		if err != nil {
			http.Error(w, errors.New("couldn't serve template").Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			Feeds           []plugin.PluginMetadata
			RefreshInterval int
		}{
			Feeds:           apiService.ListFeeds(),
			RefreshInterval: cfg.RefreshInterval,
		}
		if err := t.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(assets)))
	r.Handle("/metrics", promhttp.Handler())

	// TODO(dewey): Switch to promhttp middleware instead of this deprecated one
	r.Mount("/feed", api.NewHandler(*apiService))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("nothing to see here"))
	})

	l.Log("msg", fmt.Sprintf("feedbridge listening on http://localhost:%d", cfg.Port), "env", cfg.Environment)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
	if err != nil {
		panic(err)
	}
}
