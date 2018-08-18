package runner

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dewey/feedbridge/plugin"
	"github.com/dewey/feedbridge/scrape"
	"github.com/dewey/feedbridge/store"
	"github.com/gorilla/feeds"
)

// Runner runs scrapes
type Runner struct {
	l                    log.Logger
	Client               *http.Client
	PluginRepository     *plugin.MemRepo
	StorageRepository    *store.MemRepo
	CheckIntervalMinutes int
	ticker               *time.Ticker
}

// NewRunner initializes a new runner to run plugins
func NewRunner(l log.Logger, pluginRepo *plugin.MemRepo, storageRepo *store.MemRepo, checkIntervalMinutes int) *Runner {
	return &Runner{
		l:                    l,
		PluginRepository:     pluginRepo,
		StorageRepository:    storageRepo,
		CheckIntervalMinutes: checkIntervalMinutes,
	}
}

func init() {
	prometheus.MustRegister(scrapesDurationHistogram, pluginItemsScraped)
}

var scrapesDurationHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "fb",
		Subsystem: "runner",
		Name:      "scrape_duration_seconds",
		Help:      "Scrape duration distribution of a plugin.",
		Buckets:   []float64{5, 10, 20, 60, 120},
	},
	[]string{"plugin"},
)

var pluginItemsScraped = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "fb",
		Subsystem: "runner",
		Name:      "items_scraped",
		Help:      "Number of items scraped by a plugin.",
	},
	[]string{"plugin"},
)

// Start schedules new scrapes and runs them
func (r *Runner) Start() {
	level.Info(r.l).Log("msg", "starting plugin runner")

	if r.ticker != nil {
		r.ticker.Stop()
	}
	r.ticker = time.NewTicker(time.Duration(r.CheckIntervalMinutes) * time.Minute)
	for range r.ticker.C {
		time.Sleep(time.Duration(rand.Int63n(5000)) * time.Millisecond)
		var wg sync.WaitGroup
		for _, cp := range r.PluginRepository.All() {
			wg.Add(1)
			go func(cp plugin.Plugin) {
				start := time.Now()
				level.Info(log.With(r.l, "plugin", cp.String())).Log("msg", "scrape started")
				// TODO(dewey): Do I need to wg.Done() if it errors?
				ss, err := r.runPlugin(cp)
				if err != nil {
					level.Error(r.l).Log("err", err)
					return
				}
				wg.Done()

				duration := time.Since(start)
				scrapesDurationHistogram.WithLabelValues(cp.String()).Observe(duration.Seconds())
				pluginItemsScraped.WithLabelValues(cp.String()).Set(float64(ss.Items))
				level.Info(log.With(r.l, "plugin", cp.String())).Log("msg", "scrape finished", "feed_items", ss.Items)
			}(cp)
		}
		wg.Wait()
	}
}

func (r *Runner) runPlugin(cp plugin.Plugin) (scrape.Statistic, error) {
	f, err := cp.Run()
	if err != nil {
		return scrape.Statistic{}, err
	}
	f.Sort(func(a, b *feeds.Item) bool {
		return a.Created.After(b.Created)
	})
	rss, err := f.ToRss()
	if err != nil {
		return scrape.Statistic{}, err
	}
	r.StorageRepository.Save(fmt.Sprintf("rss_%s", cp.String()), rss)

	atom, err := f.ToAtom()
	if err != nil {
		return scrape.Statistic{}, err
	}
	r.StorageRepository.Save(fmt.Sprintf("atom_%s", cp.String()), atom)

	json, err := f.ToJSON()
	if err != nil {
		return scrape.Statistic{}, err
	}
	r.StorageRepository.Save(fmt.Sprintf("json_%s", cp.String()), json)
	return scrape.Statistic{Items: len(f.Items)}, nil
}
