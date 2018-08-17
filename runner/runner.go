package runner

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/dewey/feedbridge/repository"
	"github.com/dewey/feedbridge/store"
)

type Runner struct {
	l                    log.Logger
	Client               *http.Client
	PluginRepository     *repository.MemRepo
	StorageRepository    *store.MemRepo
	CheckIntervalMinutes int
	ticker               *time.Ticker
}

// NewRunner initializes a new runner to run plugins
func NewRunner(l log.Logger, pluginRepo *repository.MemRepo, storageRepo *store.MemRepo, checkIntervalMinutes int) *Runner {
	return &Runner{
		l:                    l,
		PluginRepository:     pluginRepo,
		StorageRepository:    storageRepo,
		CheckIntervalMinutes: checkIntervalMinutes,
	}
}

func (r *Runner) Start() {
	if r.ticker != nil {
		r.ticker.Stop()
	}
	r.ticker = time.NewTicker(time.Duration(r.CheckIntervalMinutes) * time.Second)
	for range r.ticker.C {
		time.Sleep(time.Duration(rand.Int63n(5000)) * time.Millisecond)

		var wg sync.WaitGroup
		for _, cp := range r.PluginRepository.All() {
			wg.Add(1)
			go func(cp repository.Plugin) {
				r.runPlugin(cp)
				wg.Done()
				r.l.Log("msg", "plugin finished")
			}(cp)
		}

		wg.Wait()
	}
}

func (r *Runner) runPlugin(cp repository.Plugin) {
	f, err := cp.Run()
	if err != nil {
		fmt.Println(err)
	}
	rss, err := f.ToRss()
	if err == nil {
		r.StorageRepository.Save(fmt.Sprintf("rss_%s", cp.String()), rss)
	}
	atom, err := f.ToAtom()
	if err == nil {
		r.StorageRepository.Save(fmt.Sprintf("atom_%s", cp.String()), atom)
	}
	json, err := f.ToJSON()
	if err == nil {
		r.StorageRepository.Save(fmt.Sprintf("atom_%s", cp.String()), json)
	}
}
