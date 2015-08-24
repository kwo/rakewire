# Rakewire

## v0.0.2

 - #TODO:10 move test out of src
 - #TODO:40 better logging
 - #TODO:50 add cache-control headers to API
 - #TODO:30 save feed items/entries
 - #DOING:10 add feed header to feed log (perhaps divide into http and feed subsections)
 - #TODO:60 need save database channel
 - #DONE:50 v0.0.2 Conditional GETs
 - #DONE:60 v0.0.2 feed log history - which feeds have errors, were redirected?
 - #DONE:70 v0.0.2 follow redirects, update url automatically
 - #DONE:80 v0.0.2 improve polling algorithm (automatic mode with backoff algorithm)
 - #TODO:70 bugfix feeds next on empty db or db with feeds but before first poll
 - #TODO:20 handle feeds without dates on entries - save items, check for new ones, assign new current time, update feed.LastUpdated
 - #TODO:0 bring back checksum function for comparing entries during save

## BACKLOG

 - #BACKLOG:0 have fetchers signal when they are idle, FetcherService signal when it is idle
 - #BACKLOG:10 implememt Fever API
 - #BACKLOG:20 add support for Hot Links to Fever API support
 - #BACKLOG:0 html ui
 - #BACKLOG:5 embed html ui with https://github.com/GeertJohan/go.rice
 - #BACKLOG:30 hub subscription support
 - #BACKLOG:40 TLS
 - #BACKLOG:50 http2 server push
 - #BACKLOG:60 Smart feeds based on filters (by keyword, category, author)
 - #BACKLOG:70 subscribe to web pages without a feed
 - #BACKLOG:80 plugin API for plugins for others services (Slack, XMPP, HipChat)
 - #BACKLOG:90 Twitter plugin
 - #BACKLOG:100 republish feeds (mark for publication)
 - #BACKLOG:110 hub publish support
 - #BACKLOG:120 mashup feeds, assign individual posts to new feeds
 - #BACKLOG:130 published feeds send PubSubHubbub ping
 - #BACKLOG:140 annotate entries (edit posts)
 - #BACKLOG:150 ReadLater - suck down web page and add to feed
 - #BACKLOG:160 multi-user support
 - #BACKLOG:170 automatic certificate via Let's Encrypt API
 - #BACKLOG:180 monitor mailing lists (plugin)
 - #BACKLOG:190 Full text search (Bleve)
 - #BACKLOG:200 Feeds based on filters or full-text search
 - #BACKLOG:210 database backup to tar file or something similar
 - #BACKLOG:220 limit listener to x connections
 - #BACKLOG:230 use alternative libC implementation
 - #BACKLOG:240 Mac Taskbar icon
 - #BACKLOG:250 add login for single user
 - #BACKLOG:260 pre gzip webapp
 - #BACKLOG:270 cryptographically signing of entries, feeds with nice UI lock symbol as proof (see [Atom Digital Signatures](https://tools.ietf.org/html/rfc4287#section-5.1))
 - #BACKLOG:280 manually specify feed fetchTime interval for strange feeds: Eilmeldungen, z.B.

## Changelog

 - #DONE:0 v0.0.1 config
 - #DONE:10 v0.0.1 database
 - #DONE:20 v0.0.1 httpd
 - #DONE:30 v0.0.1 fetcher pool
 - #DONE:40 v0.0.1 fetcher timer
