package runner

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/dewey/feedbridge/repository"
	"github.com/dewey/feedbridge/store"
)

type Runner struct {
	Repository           *repository.MemRepo
	Storage              *store.MemRepo
	CheckIntervalMinutes int
	ticker               *time.Ticker
	l                    log.Logger
}

// NewCountryCheck initializes a new dashboard exporter
func NewRunner(l log.Logger, checkIntervalMinutes int, repo *repository.MemRepo, storageRepo *store.MemRepo) *Runner {
	return &Runner{
		l:                    l,
		CheckIntervalMinutes: checkIntervalMinutes,
		Repository:           repo,
		Storage:              storageRepo,
	}
}

func (c *Runner) Start() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.ticker = time.NewTicker(time.Duration(c.CheckIntervalMinutes) * time.Second)
	for range c.ticker.C {
		time.Sleep(time.Duration(rand.Int63n(5000)) * time.Millisecond)

		var wg sync.WaitGroup

		for _, cp := range c.Repository.All() {
			wg.Add(1)
			go func(cp repository.Plugin) {
				c.runPlugin(cp)
				wg.Done()
				c.l.Log("msg", "plugin finished")
			}(cp)
		}

		wg.Wait()
	}
}

func (c *Runner) runPlugin(cp repository.Plugin) {
	f, err := cp.Run()
	if err != nil {
		fmt.Println(err)
	}
	rss, err := f.ToRss()
	if err == nil {
		c.Storage.Save(fmt.Sprintf("rss_%s", cp.String()), rss)
	}
	atom, err := f.ToAtom()
	if err == nil {
		c.Storage.Save(fmt.Sprintf("atom_%s", cp.String()), atom)
	}
	json, err := f.ToJSON()
	if err == nil {
		c.Storage.Save(fmt.Sprintf("atom_%s", cp.String()), json)
	}
}
