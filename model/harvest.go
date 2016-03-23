package model

// Harvest represents a polled feed complete with items and transmission data.
type Harvest struct {
	Feed         *Feed
	Items        Items
	Transmission *Transmission
}

// AddItem appends a new item to the item collection
func (z *Harvest) AddItem(guid string) *Item {
	item := I.New(z.Feed.ID, guid)
	z.Items = append(z.Items, item)
	return item
}

// Harvests is a collection of Harvest objects.
type Harvests []*Harvest
