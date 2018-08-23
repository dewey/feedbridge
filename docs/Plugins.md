# Adding a new plugin

The best starting point to get into plugin development is to look at the existing plugins
located in `plugins/`. The main idea is that a plugin gets a *http.Client to do requests
and returns a *feeds.Feed object that will then be stored by the system.

Every plugin has to implement this interface:

```
type Plugin interface {
	Run() (*feeds.Feed, error)
	Info() PluginMetadata
}
```

What you do within `Run()` is pretty flexible. You either talk to an API, Marshall JSON or you
do webscraping with `goquery`. For that usecase there are already some convinience wrappers
available, located in the `scrape` package. If sensible feel free to add new helper functions
there if they are re-usable.

# Installing the plugin

To install a plugin it has to be added to the global plugin repository, once this is done the
plugin will run periodically as defined by the interval that is set for the instance.

```
pluginRepo.Install(scmp.NewPlugin(l, c))
```

# Open a Pull Request

If you think your plugin could be useful to more people please open a Pull Request on Github:

https://github.com/dewey/feedbridge/pulls

I'll review, merge and release a new version so other people can use it on the hosted version.

