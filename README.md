# Rakewire

## v0.0.3

 - #DOING:0 html ui
 - #DOING:30 switch between dev and prod modes - version stamp files
 - #DOING:40 [react-mdl](https://github.com/tleunen/react-mdl)
 - #TODO:10 better logging
 - #TODO:20 need save database mutex
 - #FIXME:0 bugfix feeds next on empty db or db with feeds but before first poll
 - #FIXME:30 handle feeds without dates on entries - save items, check for new ones, assign new current time, update feed.LastUpdated

## BACKLOG

 - #BACKLOG:0 have fetchers signal when they are idle, FetcherService signal when it is idle
 - #BACKLOG:10 save feed items/entries
 - #BACKLOG:30 add cache-control headers to API
 - #BACKLOG:30 add cache-control headers to API
 - #BACKLOG:40 implememt Fever API
 - #BACKLOG:50 add support for Hot Links to Fever API support
 - #BACKLOG:60 hub subscription support
 - #BACKLOG:80 http2 server push
 - #BACKLOG:90 Smart feeds based on filters (by keyword, category, author)
 - #BACKLOG:100 subscribe to web pages without a feed
 - #BACKLOG:110 plugin API for plugins for others services (Slack, XMPP, HipChat)
 - #BACKLOG:120 Twitter plugin
 - #BACKLOG:130 republish feeds (mark for publication)
 - #BACKLOG:140 hub publish support
 - #BACKLOG:150 mashup feeds, assign individual posts to new feeds
 - #BACKLOG:160 published feeds send PubSubHubbub ping
 - #BACKLOG:170 annotate entries (edit posts)
 - #BACKLOG:180 ReadLater - suck down web page and add to feed
 - #BACKLOG:190 multi-user support
 - #BACKLOG:200 automatic certificate via Let's Encrypt API
 - #BACKLOG:210 monitor mailing lists (plugin)
 - #BACKLOG:220 Full text search (Bleve)
 - #BACKLOG:230 Feeds based on filters or full-text search
 - #BACKLOG:240 database backup to tar file or something similar
 - #FIXME:60 limit listener to x connections
 - #FIXME:70 use alternative libC implementation
 - #BACKLOG:270 Mac Taskbar icon
 - #BACKLOG:280 add login for single user
 - #BACKLOG:290 pre gzip webapp
 - #BACKLOG:300 cryptographically signing of entries, feeds with nice UI lock symbol as proof (see [Atom Digital Signatures](https://tools.ietf.org/html/rfc4287#section-5.1))
 - #BACKLOG:310 manually specify feed fetchTime interval for strange feeds: Eilmeldungen, z.B.
 - #BACKLOG:320 rackt redux (flux like state replay)
 - #BACKLOG:330 stamp windows executables https://github.com/josephspurrier/goversioninfo

## Changelog

 - #DONE:0 v0.0.1 config
 - #DONE:10 v0.0.1 database
 - #DONE:20 v0.0.1 httpd
 - #DONE:30 v0.0.1 fetcher pool
 - #DONE:40 v0.0.1 fetcher timer
 - #DONE:50 v0.0.2 Conditional GETs
 - #DONE:60 v0.0.2 feed log history - which feeds have errors, were redirected?
 - #DONE:70 v0.0.2 follow redirects, update url automatically
 - #DONE:80 v0.0.2 improve polling algorithm (automatic mode with backoff algorithm)
 - #DONE:90 v0.0.2 log fetches with http and feed information

## Articles

  - [Generate RESTful API Documentation From Annotations in Go](https://engineroom.teamwork.com/generate-api-from-annotations-in-go/)
  - [Solving the OPTIONS performance issue with single page apps](http://www.soasta.com/blog/options-web-performance-with-single-page-applications/?utm_source=webopsweekly&utm_medium=email)
  - [Principles of designing Go APIs with channels](https://inconshreveable.com/07-08-2014/principles-of-designing-go-apis-with-channels/)
