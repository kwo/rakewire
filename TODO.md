# Rakewire TODOs

## Roadmap

### Goals
  - Save Feed Entries
  - implememt Fever API
  - add support for Hot Links to Fever API support

### Unsorted
  - rackt redux (flux like state replay)
    - http://merrickchristensen.com/articles/single-state-tree.html
    - http://teropa.info/blog/2015/09/10/full-stack-redux-tutorial.html
    - http://konkle.us/state-management-with-redux/
    - https://github.com/rackt/redux#simple-examples
  - add client sort to feed list
  - list fetcher activity
  - have fetchers signal when they are idle, FetcherService signal when it is idle
  - Switch to Semantic-UI (like Gogs)
  - manually specify feed fetchTime interval for strange feeds: Eilmeldungen, z.B.
  - Dashboard
    - Rename home screen to Dashboard. Include various statistics including:
    - number of feeds
    - size of database
    - oldest posts

### Push
  - Add support for Sever Sent Events
    - Use SSE (less firewall problems than websockets)
    - Consider [JanBerktold/sse](https://github.com/JanBerktold/sse) and [manucorporat/sse](https://github.com/manucorporat/sse).


  - Capture Hub Feed element
	  - Capture the Hub (Atom and RSS) feed element and store it in the feedlog. Display it on the Feed info page in the UI.
  - hub subscription support

### Editing
  - add login for single user
  - Edit Feed Title and Notes fields
  - Tagging of feeds and entries

		Add support for the grouping of feeds by a user-assigned tag. Additionally, allow the tagging of individual entries.

		Starred view shows entries by tag. Unread and All display feeds by tag.


  - Add starred and read status to entries
    - Assign the flag: read and starred to entries by user.

### PushNG
  - HTTPS support
  - http2 (wait for Go 1.6)
  - See if HTTP/2 can replace Server Sent Events

### Search
  - Full text search (Bleve)
  - Feeds based on filters or full-text search
  - Smart feeds based on filters (by keyword, category, author)

### Scraping
  - Fulltext feeds - untruncate feeds
  - subscribe to web pages without a feed
  - ReadLater - suck down web page and add to feed

### Admin Console
  - multi-user support
  - database import/export
  - automatic certificate via Let's Encrypt API

### Plugins
  - plugin API for plugins for others services (Slack, XMPP, HipChat)
  - monitor mailing lists (plugin)
  - Twitter plugin

### Mashups
  - hub publish support
  - republish feeds (mark for publication)
  - mashup feeds, assign individual posts to new feeds
  - published feeds send PubSubHubbub ping
  - annotate entries (edit posts)
  - cryptographically signing of entries, feeds with nice UI lock symbol as proof
    - see [Atom Digital Signatures](https://tools.ietf.org/html/rfc4287#section-5.1)

### Installation
  - Make rakewire run as service
    - Consider https://github.com/kardianos/service
  - Mac Taskbar icon

### Maintenance
  - Limit listener to x connections
  - Review channels in app
    - See [principles of designing go apis with channels](https://inconshreveable.com/07-08-2014/principles-of-designing-go-apis-with-channels/)
  - Ensure no problems with OPTIONS
    - See [Options Web Performance With SPAs](http://www.soasta.com/blog/options-web-performance-with-single-page-applications/)
  - Feeds view: slim size of data payload
  - consider musl as alternative libC implementation
    - consider [Musl](Consider Musl: http://www.musl-libc.org/)
    - See [article](http://dominik.honnef.co/posts/2015/06/statically_compiled_go_programs__always__even_with_cgo__using_musl/)

### DevOps
  - Git Repository
    - create docker images
    - document upgrade procedure
    - move Gogs to own server
  - CI to test binary
    - See [six continuous integration tools](http://opensource.com/business/15/7/six-continuous-integration-tools)