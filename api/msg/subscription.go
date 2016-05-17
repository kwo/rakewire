package msg

import (
	"time"
)

// Subscriptions is a list of Subscription structs
type Subscriptions []*Subscription

// Subscription defines a subscription to a Feed
type Subscription struct {
	URL      string    `json:"url,omitempty"`
	Title    string    `json:"title,omitempty"`
	Groups   []string  `json:"groups,omitempty"`
	Notes    string    `json:"notes,omitempty"`
	Added    time.Time `json:"added,omitempty"`
	AutoRead bool      `json:"autoread,omitempty"`
	AutoStar bool      `json:"autostar,omitempty"`
}

// SubscriptionAddUpdateRequest defines an add/update subscription request
type SubscriptionAddUpdateRequest struct {
	AddGroups    bool          `json:"addGroups"`
	Subscription *Subscription `json:"subscription"`
}

// SubscriptionAddUpdateResponse defines the response to a SubscriptionAddUpdateRequest
type SubscriptionAddUpdateResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

// SubscriptionListRequest defines the request to add a subscription
type SubscriptionListRequest struct {
	Filter string `json:"filter,omitempty"`
}

// SubscriptionListResponse returns a list of subscriptions
type SubscriptionListResponse struct {
	Status        int           `json:"status"`
	Message       string        `json:"message,omitempty"`
	Subscriptions Subscriptions `json:"subscriptions,omitempty"`
}

// UnsubscribeRequest defines the request to remove a subscription
type UnsubscribeRequest struct {
	URL string `json:"url,omitempty"`
}

// UnsubscribeResponse defines the response to a UnsubscribeRequest
type UnsubscribeResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}
