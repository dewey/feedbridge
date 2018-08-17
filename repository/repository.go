package repository

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
	String() string
}
