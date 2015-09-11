# Rakewire

## v0.0.3

 - TODO: remove negroni, [use normal http handlers](https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81)
 - TODO: html ui
 - HACK: load roboto font locally
 - TODO: set cache-control headers on static assets
 - XXX: list feed log
 - XXX: add client sort to feed list
 - TODO: list fetcher activity
 - TODO: tests
 - TODO: better logging
 - BACKLOG: http2
 - HACK: need save database mutex
 - TODO: bugfix feeds next on empty db or db with feeds but before first poll
 - TODO: handle feeds without dates on entries - save items, check for new ones, assign new current time, update feed.LastUpdated

## IDEAS

 - BACKLOG: have fetchers signal when they are idle, FetcherService signal when it is idle
 - BACKLOG: save feed items/entries
 - BACKLOG: add cache-control headers to API
 - BACKLOG: add cache-control headers to API
 - BACKLOG: implememt Fever API
 - BACKLOG: add support for Hot Links to Fever API support
 - BACKLOG: hub subscription support
 - BACKLOG: http2 server push or SSE
 - BACKLOG: Smart feeds based on filters (by keyword, category, author)
 - BACKLOG: subscribe to web pages without a feed
 - BACKLOG: plugin API for plugins for others services (Slack, XMPP, HipChat)
 - BACKLOG: Twitter plugin
 - BACKLOG: republish feeds (mark for publication)
 - BACKLOG: hub publish support
 - BACKLOG: mashup feeds, assign individual posts to new feeds
 - BACKLOG: published feeds send PubSubHubbub ping
 - BACKLOG: annotate entries (edit posts)
 - BACKLOG: ReadLater - suck down web page and add to feed
 - BACKLOG: multi-user support
 - BACKLOG: automatic certificate via Let's Encrypt API
 - BACKLOG: monitor mailing lists (plugin)
 - BACKLOG: Full text search (Bleve)
 - BACKLOG: Feeds based on filters or full-text search
 - BACKLOG: database backup to tar file or something similar
 - BACKLOG: limit listener to x connections
 - BACKLOG: use alternative libC implementation
 - BACKLOG: Mac Taskbar icon
 - BACKLOG: add login for single user
 - BACKLOG: pre gzip webapp
 - BACKLOG: cryptographically signing of entries, feeds with nice UI lock symbol as proof (see [Atom Digital Signatures](https://tools.ietf.org/html/rfc4287#section-5.1))
 - BACKLOG: manually specify feed fetchTime interval for strange feeds: Eilmeldungen, z.B.
 - BACKLOG: rackt redux (flux like state replay)
 - BACKLOG: stamp windows executables https://github.com/josephspurrier/goversioninfo

## Changelog

 - DONE: v0.0.1 config
 - DONE: v0.0.1 database
 - DONE: v0.0.1 httpd
 - DONE: v0.0.1 fetcher pool
 - DONE: v0.0.1 fetcher timer
 - DONE: v0.0.2 Conditional GETs
 - DONE: v0.0.2 feed log history - which feeds have errors, were redirected?
 - DONE: v0.0.2 follow redirects, update url automatically
 - DONE: v0.0.2 improve polling algorithm (automatic mode with backoff algorithm)
 - DONE: v0.0.2 log fetches with http and feed information
