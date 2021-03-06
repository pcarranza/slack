package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// StartRTM calls the "rtm.start" endpoint and returns the provided URL and the full Info block.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) StartRTM() (info *Info, websocketURL string, err error) {
	return api.StartRTMContext(context.Background())
}

// StartRTMContext calls the "rtm.start" endpoint and returns the provided URL and the full Info block with a custom context.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) StartRTMContext(ctx context.Context) (info *Info, websocketURL string, err error) {
	response := &infoResponseFull{}
	err = post(ctx, "rtm.start", url.Values{"token": {api.config.token}}, response, api.debug)
	if err != nil {
		return nil, "", fmt.Errorf("post: %s", err)
	}
	if !response.Ok {
		return nil, "", response.Error
	}
	api.Debugln("Using URL:", response.Info.URL)
	return &response.Info, response.Info.URL, nil
}

// ConnectRTM calls the "rtm.connect" endpoint and returns the provided URL and the compact Info block.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) ConnectRTM() (info *Info, websocketURL string, err error) {
	return api.ConnectRTMContext(context.Background())
}

// ConnectRTM calls the "rtm.connect" endpoint and returns the provided URL and the compact Info block with a custom context.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()` on it.
func (api *Client) ConnectRTMContext(ctx context.Context) (info *Info, websocketURL string, err error) {
	response := &infoResponseFull{}
	err = post(ctx, "rtm.connect", url.Values{"token": {api.config.token}}, response, api.debug)
	if err != nil {
		return nil, "", fmt.Errorf("post: %s", err)
	}
	if !response.Ok {
		return nil, "", response.Error
	}
	api.Debugln("Using URL:", response.Info.URL)
	return &response.Info, response.Info.URL, nil
}

// NewRTM returns a RTM, which provides a fully managed connection to
// Slack's websocket-based Real-Time Messaging protocol.
func (api *Client) NewRTM() *RTM {
	return api.NewRTMWithOptions(nil)
}

// NewRTMWithOptions returns a RTM, which provides a fully managed connection to
// Slack's websocket-based Real-Time Messaging protocol.
// This also allows to configure various options available for RTM API.
func (api *Client) NewRTMWithOptions(options *RTMOptions) *RTM {
	result := &RTM{
		Client:           *api,
		IncomingEvents:   make(chan RTMEvent, 50),
		outgoingMessages: make(chan OutgoingMessage, 20),
		pings:            make(map[int]time.Time),
		isConnected:      false,
		wasIntentional:   true,
		killChannel:      make(chan bool),
		forcePing:        make(chan bool),
		rawEvents:        make(chan json.RawMessage),
		idGen:            NewSafeID(1),
	}

	if options != nil {
		result.useRTMStart = options.UseRTMStart
	} else {
		result.useRTMStart = true
	}

	return result
}
