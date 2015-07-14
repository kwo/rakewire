# Rakewire

## TODO
 - implememt APIRouter to list feeds in database
 - run fetch off of database
 - get fetch to complete without errors
 - start fetcher pool, feed fetcher in addition to database and httpd
 - implement config defaults

## Version 1 (MVP)

 - Feed Reading via Browser
 - implememt Fever API

## Roadmap

 - Conditional GETs
 - Set Feed fetch frequency (automatic mode with backoff algorithm)
 - Push to Browser via WebSockets, etc.
 - PubSubHubBub
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
 - Mac Taskbar icon

## Technical Don't Forgets
 - use alternative libC implementation
 - limit listener to x connections
