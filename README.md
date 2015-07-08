# Rakewire

## TODO
 - open database at start
 - add db location to config file
 - decide on a database location
 - implement config defaults
 - implememt APIRouter to list feeds in database
 - trap signal
 - close database on signal
 - gracefully stop server on signal
 - run fetch off of database
 - get fetch to complete without errors

## Version 1 (MVP)

 - Feed Reading via Browser
 - implememt Fever API

## Roadmap

 - Conditional GETs
 - Set Feed fetch frequency (automatic mode with backoff algorithm)
 - PubSubHubBub
 - Push to Browser via WebSockets, etc.
 - Feed filters (by keyword, category, author)
 - Full text search (Bleve)
 - Feeds based on filters or full-text search
 - subscribe to web pages without a feed
 - read Twitter feeds
 - monitor mailing lists
 - plugins for others services (Slack, XMPP, HipChat)
 - ReadLater - suck down web page
 - Publish feeds
 - TLS via Let's Encrypt
 - HTTP2
 - database backup to tar file or something similar

## Technical Don't Forgets
 - use alternative libC implementation
