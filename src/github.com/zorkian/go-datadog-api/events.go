/*
 * Datadog API for Go
 *
 * Please see the included LICENSE file for licensing information.
 *
 * Copyright 2013 by authors and contributors.
 */

package datadog

import (
	"fmt"
	"net/url"
	"strconv"
)

// Event is a single event. If this is being used to post an event, then not
// all fields will be filled out.
type Event struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	Text        string   `json:"text"`
	Time        int      `json:"date_happened"` // UNIX time.
	Priority    string   `json:"priority"`
	AlertType   string   `json:"alert_type"`
	Host        string   `json:"host"`
	Aggregation string   `json:"aggregation_key"`
	SourceType  string   `json:"source_type"`
	Tags        []string `json:"tags"`
	Url         string   `json:"url"`
	Resource    string   `json:"resource"`
}

// reqGetEvent is the container for receiving a single event.
type reqGetEvent struct {
	Event Event `json:"event"`
}

// reqGetEvents is for returning many events.
type reqGetEvents struct {
	Events []Event `json:"events"`
}

// PostEvent takes as input an event and then posts it to the server.
func (self *Client) PostEvent(event *Event) (*Event, error) {
	var out reqGetEvent
	err := self.doJsonRequest("POST", "/v1/events", event, &out)
	if err != nil {
		return nil, err
	}
	return &out.Event, nil
}

// GetEvent gets a single event given an identifier.
func (self *Client) GetEvent(id int) (*Event, error) {
	var out reqGetEvent
	err := self.doJsonRequest("GET", fmt.Sprintf("/v1/events/%d", id), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out.Event, nil
}

// QueryEvents returns a slice of events from the query stream.
func (self *Client) GetEvents(start, end int,
	priority, sources, tags string) ([]Event, error) {
	// Since this is a GET request, we need to build a query string.
	vals := url.Values{}
	vals.Add("start", strconv.Itoa(start))
	vals.Add("end", strconv.Itoa(end))
	if priority != "" {
		vals.Add("priority", priority)
	}
	if sources != "" {
		vals.Add("sources", sources)
	}
	if tags != "" {
		vals.Add("tags", tags)
	}

	// Now the request and response.
	var out reqGetEvents
	err := self.doJsonRequest("GET",
		fmt.Sprintf("/v1/events?%s", vals.Encode()), nil, &out)
	if err != nil {
		return nil, err
	}
	return out.Events, nil
}
