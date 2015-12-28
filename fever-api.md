# Fever API

## root url

	?api[=xml]

## Read

### groups (DONE)

	?api&groups

Response contains groups and feeds_groups.

A group object has the following members:
id (positive integer)
title (utf-8 string)

A feeds_group object has the following members:
group_id (positive integer)
feed_ids (string/comma-separated list of positive integers)


### feeds (DONE)

	?api&feeds

Response contains feeds and feeds_groups.

A feeds_group object has the following members:
group_id (positive integer)
feed_ids (string/comma-separated list of positive integers)

A feed object has the following members:
id (positive integer)
favicon_id (positive integer)
title (utf-8 string)
url (utf-8 string)
site_url (utf-8 string)
is_spark (boolean integer)
last_updated_on_time (Unix timestamp/integer)

### favicons

	?api&favicons

Response contains favicons.

### items (DONE)

	?api&items [&since_id=] [&max_id=] [&with_ids=id,id]

Response contains items and total_items.


Use the since_id argument with the highest id of locally cached items to request 50 additional items. Repeat until the items array in the response is empty.

Use the max_id argument with the lowest id of locally cached items (or 0 initially) to request 50 previous items. Repeat until the items array in the response is empty. (added in API version 2)

Use the with_ids argument with a comma-separated list of item ids to request (a maximum of 50) specific items. (added in API version 2)


total_items: contains the total number of items stored in the database

An item object has the following members:
id (positive integer)
feed_id (positive integer)
title (utf-8 string)
author (utf-8 string)
html (utf-8 string)
url (utf-8 string)
is_saved (boolean integer)
is_read (boolean integer)
created_on_time (Unix timestamp/integer)


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
