package racefansnet

import (
	"net/http"

	pm "github.com/dewey/feedbridge/plugin"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

// Plugin defines a new plugin
type plugin struct {
	l log.Logger
	c *http.Client
	f *feeds.Feed
}

// NewPlugin initializes a new plugin
func NewPlugin(l log.Logger, c *http.Client) *plugin {
	return &plugin{
		l: log.With(l, "plugin", "racefansnet"),
		c: c,
		f: &feeds.Feed{
			Title:       "Racefans.net",
			Link:        &feeds.Link{Href: "https://www.racefans.net"},
			Description: "Independend motorsport coverage.",
			Author:      &feeds.Author{Name: "Facefans.net Authors"},
		},
	}
}

func (p *plugin) Info() pm.PluginMetadata {
	return pm.PluginMetadata{
		TechnicalName: "racefansnet",
		Name:          p.f.Title,
		Description:   "Providing a working feed, WP feed is broken.",
		Author:        "Philipp",
		AuthorURL:     "https://github.com/dewey",
		Image:         "https://i.imgur.com/xYTKND4.png",
		SourceURL:     "https://www.racefans.net",
	}
}

func (p *plugin) Run() (*feeds.Feed, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.racefans.net/feed/", nil)
	if err != nil {
		return nil, err
	}
	// If these headers are not set the feed doesn't load properly, thanks to Wordpress probably
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:67.0) Gecko/20100101 Firefox/67.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var feedItems []*feeds.Item
	for _, fi := range feed.Items {
		item := &feeds.Item{
			Author: &feeds.Author{
				Name:  fi.Author.Name,
				Email: fi.Author.Email,
			},
			Title: fi.Title,
			Link: &feeds.Link{
				Href: fi.Link,
			},
			Id:          fi.GUID,
			Description: fi.Description,
		}
		if fi.PublishedParsed != nil {
			item.Created = *fi.PublishedParsed
		}
		if fi.UpdatedParsed != nil {
			item.Updated = *fi.UpdatedParsed
		}
		feedItems = append(feedItems, item)
	}
	p.f.Items = feedItems
	return p.f, nil
}
