package scmp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

// Plugin defines a new plugin
type plugin struct {
	l log.Logger
}

// NewChecker initializes a new dashboard exporter
func NewPlugin(l log.Logger) *plugin {
	return &plugin{
		l: l,
	}
}

func (p *plugin) String() string {
	return "scmp"
}

// Run runs the main checker function of the plugin
func (p *plugin) Run() (*feeds.Feed, error) {
	c := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", "https://www.scmp.com/topics/infographics-asia", nil)
	if err != nil {
		return nil, err
	}
	// They block requests without valid user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.1.2 Safari/605.1.15")
	req.Header.Set("Accept-Language", "en-us")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err

	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	feed := &feeds.Feed{
		Title:       "South China Morning Post",
		Link:        &feeds.Link{Href: "https://www.scmp.com/topics/infographics-asia"},
		Description: "Your source for credible news and authoritative insights from Hong Kong, China and the world.",
		Author:      &feeds.Author{Name: "SCMP", Email: "digitalsupport@scmp.com"},
	}

	var feedItems []*feeds.Item
	doc.Find("div.pane-article-level article").Each(func(i int, s *goquery.Selection) {
		item := &feeds.Item{
			Author: feed.Author,
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
				fmt.Println(err)
			}
		}

		feedItems = append(feedItems, item)
	})
	feed.Items = feedItems
	return feed, nil
}
