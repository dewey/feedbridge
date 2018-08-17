package scmp

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dewey/feedbridge/scrape"
	"github.com/go-kit/kit/log"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

// Plugin defines a new plugin
type plugin struct {
	l log.Logger
	c *http.Client
	f *feeds.Feed
}

// NewChecker initializes a new dashboard exporter
func NewPlugin(l log.Logger, c *http.Client) *plugin {
	return &plugin{
		l: l,
		c: c,
		f: &feeds.Feed{
			Title:       "South China Morning Post",
			Link:        &feeds.Link{Href: "https://www.scmp.com/topics/infographics-asia"},
			Description: "Your source for credible news and authoritative insights from Hong Kong, China and the world.",
			Author:      &feeds.Author{Name: "SCMP", Email: "digitalsupport@scmp.com"},
		},
	}
}

func (p *plugin) String() string {
	return "scmp"
}

// Run runs the main checker function of the plugin
func (p *plugin) Run() (*feeds.Feed, error) {
	it := []string{
		"https://www.scmp.com/topics/infographics-asia",
		"https://www.scmp.com/topics/infographics-politics",
		"https://www.scmp.com/topics/infographics-lifestyle",
		"https://www.scmp.com/topics/infographics-international",
		"https://www.scmp.com/topics/infographics-science",
		"https://www.scmp.com/topics/infographics-economics",
	}

	result, err := scrape.URLToDocument(p.c, scrape.URLtoTask(it))
	if err != nil {
		return nil, err
	}
	var feedItems []*feeds.Item
	for _, r := range result {
		// Get all top level items
		items, err := p.listHandler(&r.Document)
		if err != nil {
			p.l.Log("err", err)
		}
		feedItems = append(feedItems, items...)

		// Create tasks for pagination
		var subTask []scrape.Task
		for i := 1; i < 2; i++ {
			u, err := url.Parse(r.URL)
			if err != nil {
				fmt.Println(err)
				continue
			}
			q := u.Query()
			q.Add("page", strconv.Itoa(i))
			u.RawQuery = q.Encode()
			subTask = append(subTask, scrape.Task{
				URL: u.String(),
			})
		}

		// Get all items from other pages
		result, err := scrape.URLToDocument(p.c, subTask)
		if err != nil {
			return nil, err
		}
		for _, r := range result {
			items, err := p.listHandler(&r.Document)
			if err != nil {
				p.l.Log("err", err)
			}
			feedItems = append(feedItems, items...)
		}
	}
	p.f.Items = feedItems
	return p.f, nil
}

func (p *plugin) listHandler(doc *goquery.Document) ([]*feeds.Item, error) {
	var feedItems []*feeds.Item
	doc.Find("div.pane-article-level article").Each(func(i int, s *goquery.Selection) {
		item := &feeds.Item{
			Author: p.f.Author,
		}

		val, exists := s.Attr("about")
		if exists {
			url := fmt.Sprintf("https://www.scmp.com%s", val)
			item.Link = &feeds.Link{Href: url}
			// Unique identifier helps the client with deduplicating
			item.Id = url
		}

		ds := s.Find("div.content-wrapper > div.caption-wrapper > h3.node-title > a")
		item.Title = ds.Text()

		is := s.Find("div.background-image > a > img")
		val, exists = is.Attr("data-original")
		if exists {
			item.Description = fmt.Sprintf(`<img src="%s">`, val)
		}

		ts := s.Find("span.rdf-meta")
		val, exists = ts.Attr("property")
		if exists {
			if val == "dc:title" {
				val, exists = ts.Attr("content")
				if exists {
					item.Description = fmt.Sprintf(`%s<p>%s`, item.Description, val)
				}
			}
		}
		times := s.Find("time.updated")
		val, exists = times.Attr("content")
		if exists {
			t, err := time.Parse("2006-01-02T15:04:05-07:00", val)
			if err == nil {
				item.Updated = t
			} else {
				p.l.Log("err", err)
			}
		}

		feedItems = append(feedItems, item)
	})
	return feedItems, nil
}
