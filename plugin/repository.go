package plugin

import (
	"github.com/gorilla/feeds"
)

// A Repository is responsible for retrieving/registering pluginFuncs
type Repository interface {
	Install(p Plugin)
	Find(name string) (Plugin, error)
	All() []Plugin
}

// Plugin is the interface that a scrape plugin has to implement
type Plugin interface {
	Run() (*feeds.Feed, error)
	Info() PluginMetadata
}

type PluginMetadata struct {
	TechnicalName string
	Name          string
	Description   string
	Author        string
	AuthorURL     string
	Image         string
	SourceURL     string
}
