/*
Copyright (C) 2018 Expedia Group.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"github.com/tbruyelle/hipchat-go/hipchat"
	"net/http"
	"time"
)

type HipchatClient interface {
	SendMessage(roomID, message string) error
	SendNotification(roomID string, notification *hipchat.NotificationRequest) error
	GetMessages(roomID string, options *hipchat.LatestHistoryOptions) ([]hipchat.Message, error)
}

type hipchatClient struct {
	clientPool chan *hipchat.Client
}

func NewHipChatClient(authTokens []string) HipchatClient {

	pool := make(chan *hipchat.Client, len(authTokens))
	for _, t := range authTokens {
		hc := hipchat.NewClient(t)
		hc.SetHTTPClient(&http.Client{
			Timeout: time.Second * 15,
		})
		pool <- hc
	}
	return hipchatClient{clientPool: pool}
}

// Always returnClient after use to make it available again
func (c *hipchatClient) getClient() *hipchat.Client {
	return <-c.clientPool
}

// Return the client for use by other operations
func (c *hipchatClient) returnClient(client *hipchat.Client) {
	// Return the token for use after a few seconds delay - limiting API calls/minute
	go func() {
		time.Sleep(5 * time.Second)
		c.clientPool <- client
	}()
}

func (c hipchatClient) SendMessage(roomID, message string) error {

	hcl := c.getClient()
	defer c.returnClient(hcl)

	messageRequest := &hipchat.RoomMessageRequest{Message: message}
	err := do(func() error {
		_, err := hcl.Room.Message(roomID, messageRequest)
		if err != nil {
			time.Sleep(5 * time.Second)
			return err
		}
		return nil
	})

	return err
}

func (c hipchatClient) GetMessages(roomID string, options *hipchat.LatestHistoryOptions) ([]hipchat.Message, error) {

	hcl := c.getClient()
	defer c.returnClient(hcl)

	if history, _, err := hcl.Room.Latest(roomID, options); err != nil {
		return []hipchat.Message{}, err
	} else {
		return history.Items, nil
	}
}

func (c hipchatClient) SendNotification(roomID string, notification *hipchat.NotificationRequest) error {

	hcl := c.getClient()
	defer c.returnClient(hcl)

	err := do(func() error {
		_, err := hcl.Room.Notification(roomID, notification)
		if err != nil {
			time.Sleep(5 * time.Second)
		}
		return nil
	})

	return err
}
