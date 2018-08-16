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

// Plugin is the interface used in each plugin (handler)
type Plugin interface {
	Run() (*feeds.Feed, error)
	String() string
}
