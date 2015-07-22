# Rakewire

## TODO

## Version 0.0.1

  - DONE config
  - DONE database
  - DONE httpd
  - DONE fetcher pool
  - DONE fetcher timer

## Version 0.0.2
  - Conditional GETs
  - feed log history - which feeds have errors, were redirected?
  - feed parser
  - save feed items/entries

## Version 0.1.0
  - rest api
  - html ui

## Version 1.0.0 (MVP)

 - implememt Fever API
 - compatiability with Reeder2

## Version 1.1.0 (MVP)

 - add support for Hot Links to Fever API support

## Version 2.0.0
 - websocket api
 - hub support

## Roadmap

 - Feed Reading via Browser (UI)
 - Set Feed fetch frequency (automatic mode with backoff algorithm)
 - Feed Feedback: update url automatically
 - Push to Browser via WebSockets, etc.
 - PubSubHubBub
 - Feed filters (by keyword, category, author)
 - Full text search (Bleve)
 - Feeds based on filters or full-text search
 - subscribe to web pages without a feed
 - Support for multiple users

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
 - need save database channel
 - use alternative libC implementation
 - limit listener to x connections
