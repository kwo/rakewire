# Rakewire Changelog

## 1.6.1 2016-01-28

 - Bugfix: add back httpd.address parameter to specify a binding address separate from the hostname

## 1.6.0 2016-01-28

 * read config from database,
   allowing for a single data file
 * add status page to /status
 * add version flag
 * force integrity check when app version changes
 * load tls certificate via config

## 1.5.0 2016-01-27

 * add -check flag to check database integrity
 * Entry links default to alternate if rel missing
 * bugfix: remove Delete from Cursor
 * remove rest cleanup method
 * create collection classes

## 1.4.0 2016-01-26

 * move data access to model package,
   require transaction for every operation  
 * rename Entry to Item
 * rename FeedLog to Transmission
 * rename UserFeed to Subscription
 * rename UserEntry to Entry

## 1.3.1 2016-01-20

 * bugfix: disabled feedparser filter, needs work

## 1.3.0 2016-01-20

 * filter control characters from feeds
 * set subscription title to url if no title available

## 1.2.0 2016-01-17

 * add basic authentication to rest API

## 1.1.0 2016-01-14

 * sort OPML export

## 1.0.2 2016-01-14

 * add gzip to fever api
 * shut off static pages and old api
 * add xml content type to opml export
 * add + symbol to opml tags

## 1.0.1 2016-01-13 LAUNCH RAKEWIRE!

 * add OPML support
 * adjust fetch time to 15 minutes for next 2 hours
 * max next fetch time now 1 hour
 * prevent entries with future dates
 * bugfix: add support for marking items as unread
 * bugfix: mark fever individual items as read
 * bugfix: feeds_groups not feed_groups

## 1.0.0 2016-01-05
