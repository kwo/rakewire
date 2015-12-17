# Fever API

## root url

	?api[=xml]

## Read

### groups (TODO)

	?api&groups

Response contains groups and feeds_groups.


### feeds (TODO)

	?api&feeds

Response contains feeds and feeds_groups.

### favicons

	?api&favicons

Response contains favicons.

### items (TODO)

	?api&items [&since_id=] [&max_id=] [&with_ids=id,id]

Response contains items and total_items.


### hot links

	?api&links

Response contains links.

### sync (TODO)

	?api&unread_item_ids
	?api&saved_item_ids


## Write

### unread

I don't understand this yet.

	?api

	unread_recently_read=1


### items (TODO)

	?api

	mark=item
	as=? where ? is replaced with read, saved or unsaved
	id=? where ? is replaced with the id of the item to modify


### feed or group (TODO)

	?api

	mark=? where ? is replaced with feed or group
	as=read
	id=? where ? is replaced with the id of the feed or group to modify
	before=? where ? is replaced with the Unix timestamp of the the local client’s most recent items API request


### kindling

	?api

	mark=group
	as=read
	id=0
	before=? where ? is replaced with the Unix timestamp of the the local client’s last items API request


### sparks

	?api

	mark=group
	as=read
	id=-1
	before=? where ? is replaced with the Unix timestamp of the the local client’s last items API request
