<img src="https://i.imgur.com/0j5LG0t.png" alt="feedbridge gopher" width="100">

# feedbridge

[![Build Status](https://travis-ci.com/dewey/feedbridge.svg?branch=master)](https://travis-ci.com/dewey/feedbridge)
[![GoDoc](https://godoc.org/github.com/dewey/feedbridge?status.svg)](https://godoc.org/github.com/dewey/feedbridge) 
[![Go Report Card](https://goreportcard.com/badge/github.com/dewey/feedbridge)](https://goreportcard.com/report/github.com/dewey/feedbridge)
[![Maintainability](https://api.codeclimate.com/v1/badges/50d72195e5d1f42d21b1/maintainability)](https://codeclimate.com/github/dewey/feedbridge/maintainability)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dewey/feedbridge/LICENSE)
[![Badges](https://img.shields.io/badge/badges-all%20of%20them-brightgreen.svg)](https://github.com/dewey/feedbridge)

Is a tool to provide RSS feeds for sites that don't have one, or only offer a feed of headlines. For each site—or kind of site—you want to generate a feed for you'll have to implement a plugin with a custom scraping strategy. feedbridge doesn't persist old items so if it's not on the site you are scraping any more it won't be in the feed. Pretty similar to how most feeds these days work that only have the latest items in there. It publishes Atom, RSS 2.0, and JSON Feed Version 1 conform feeds.

There are a bunch of web apps doing something similar, some of them you can even drag and drop selectors to create a feed. That didn't work well for the site I was trying it for so I decided to built this. (Also it was fun doing so).

## API

**GET /feed/list**

Returns a list of available plugins.

**GET /feed/{plugin}/{format}**

Returns the feed based on a given plugin and output format. That's the URL you should use in your feed reader.

- `plugin`: The name of the plugin as returned by `String()`
- `format`: The format the feed should be returned in, can be `rss`, `atom` or `json`. By default it's RSS.

**GET /metrics**

Returns the exported Prometheus metrics.

## Configuration and Operation

### Environment

The following environment variables are available, they all have sensible defaults and don't need to be set explicity.

- `REFRESH_INTERVAL`: The interval in which feeds get rescraped in minutes (Default: 15)
- `ENVIRONMENT`: Environment can be `prod` or `develop`. `develop` sets the loglevel to `info` (Default: `develop`)
- `PORT`: Port that feedbridge is running on (Default: `8080`)

There are two available storage backends right now. An in-memory and a disk backed implementation. Depending on which one you choose
there are additional options you can set.

- `STORAGE_BACKEND`: Set to `memory` to keep everything in-memory or `persistent` to persist the cache to disk. (Default: `memory`)

**In Memory**

- `CACHE_EXPIRATION`: The expiration time of the cache in minutes (Default: 30)
- `CACHE_EXPIRED_PURGE`: The interval at which the expired cache elements will be purged in minutes (Default: 60)

**Persistent**

- `STORAGE_PATH`: Set the storage location of the cache on disk. (Default: `./feedbridge-cache`)


### Run with Docker

You can change all these options in the included docker-compose files and use `docker-compose -f docker-compose.yml up -d` to run the project. There are
two sample docker compose files that already have the settings for the two available backends set.

### Monitoring

As the project already exports Prometheus metrics you can use Grafana to get more information about how many things are being scraped and how fast requests are served. You can import the included `grafana-dashboard.json` in Grafana.

## Status

This is a work in progress and pretty rough right now. The API might change and things get moved around.

## Acknowledgements & Credits

It's using the neat [gorilla/feeds](https://github.com/gorilla/feeds) library to generate standard conform Atom, RSS 2.0 and JSON Feeds. The Gopher was sourced from [github.com/egonelbre](https://github.com/egonelbre/gophers), the RSS icon is coming from Wikipedia and was added by me. Thanks!