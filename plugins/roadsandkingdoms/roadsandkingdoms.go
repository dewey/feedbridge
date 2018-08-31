package roadsandkingdoms

import (
	"io/ioutil"
	"net/http"

	pm "github.com/dewey/feedbridge/plugin"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/feeds"
	readability "github.com/mauidude/go-readability"
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
		l: log.With(l, "plugin", "roadsandkingdoms"),
		c: c,
		f: &feeds.Feed{
			Title:       "Roads & Kingdoms",
			Link:        &feeds.Link{Href: "https://roadsandkingdoms.com"},
			Description: "Journalism and travel, together at last.",
			Author:      &feeds.Author{Name: "Roads & Kingdoms"},
		},
	}
}

func (p *plugin) Info() pm.PluginMetadata {
	return pm.PluginMetadata{
		TechnicalName: "roadsandkingdoms",
		Name:          p.f.Title,
		Description:   "Proving a full content feed not just snippets.",
		Author:        "Philipp",
		AuthorURL:     "https://github.com/dewey",
		Image:         "https://i.imgur.com/ABzvg51.png",
		SourceURL:     "https://roadsandkingdoms.com",
	}
}

func (p *plugin) Run() (*feeds.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://roadsandkingdoms.com/feed/")
	if err != nil {
		return nil, err
	}

	var feedItems []*feeds.Item
	for _, fi := range feed.Items {
		resp, err := p.c.Get(fi.Link)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		doc, err := readability.NewDocument(string(b))
		if err != nil {
			continue
		}

		content := doc.Content()
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
			Description: content,
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
