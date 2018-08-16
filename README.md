# feedbridge

Is a tool to provide RSS feeds for sites that don't have one. For each site you want you'll have to create a plugin with a custom scraping strategy. feedbridge doesn't persist old items so if it's not on the site you are scraping any more it won't be in the feed. Pretty similar to how most feeds these days work that only have the latest items in there. It's using the neat [gorilla/feeds](https://github.com/gorilla/feeds) library which is supporting Atom, RSS 2.0, and JSON Feed Version 1 spec elements.

There are a bunch of web apps doing something similar, some of them you can even drag and drop selectors to create a feed. That didn't work well for the site I was trying it for so I decided to built this. (Also it was fun doing so).

## API

**GET /feed/list**

Returns a list of available plugins.

**GET /feed/{plugin}/{format}**

Returns the feed based on a given plugin and output format.

- `plugin`: The name of the plugin as returned by `String()`
- `format`: The format the feed should be returned in, can be `rss`, `atom` or `json`. By default it's RSS.