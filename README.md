# Rakewire

## TODO
 - put routes in separate file
 - implememt APIRouter to list feeds in database
 - run fetch off of database
 - start fetcher pool, feed fetcher in addition to database and httpd
 - start timer every 5 minutes to repoll feeds
 - get fetch to complete without errors
 - implement config defaults

## Version 0.0.1

  - DONE config
  - DONE database
  - DONE httpd
  - rest api
  - fetcher pool
  - fetcher timer
  - html ui

## Version 1.0.0 (MVP)

 - implememt Fever API
 - compatiability with Reeder2

## Version 1.1.0 (MVP)

 - add support for Hot Links to Fever API support

## Roadmap

 - Feed Reading via Browser (UI)
 - Conditional GETs
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
 - use alternative libC implementation
 - limit listener to x connections
