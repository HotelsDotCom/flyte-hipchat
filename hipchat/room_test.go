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

package hipchat

import (
	"github.com/stretchr/testify/assert"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"sync/atomic"
	"testing"
	"time"
)

func TestLeave(t *testing.T) {

	messagesOut := make(chan Message)
	cm := NewClientMock()
	cm.getMessages = func(string, *hipchat.LatestHistoryOptions) ([]hipchat.Message, error) {
		return []hipchat.Message{{Message: "incoming message"}}, nil
	}

	room := NewRoom("abc", cm, messagesOut)

	// event handler consuming messages
	var counter int32 = 0
	go func() {
		for {
			<-messagesOut
			atomic.AddInt32(&counter, 1)
			counter += 1
		}
	}()

	time.Sleep(10 * time.Millisecond) // process some message
	room.Leave()
	time.Sleep(10 * time.Millisecond) // wait for leave signal to be picked up

	// no more messages should be processed after leave
	currentCounter := counter
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, currentCounter, counter, "did not expect to process any more messages after leave")
	assert.NotEqual(t, int32(0), counter, "no messages processed")
}

// --- mocks ---

type SendMessageCall struct {
	roomId  string
	message string
}

type SendNotificationCall struct {
	roomId       string
	notification *hipchat.NotificationRequest
}

type GetMessagesCall struct {
	roomID  string
	options *hipchat.LatestHistoryOptions
}

type ClientMock struct {
	SendMessageCall      SendMessageCall
	SendNotificationCall SendNotificationCall
	GetMessagesCall      GetMessagesCall

	sendMessage      func(roomID, message string) error
	sendNotification func(roomID string, notification *hipchat.NotificationRequest) error
	getMessages      func(roomID string, options *hipchat.LatestHistoryOptions) ([]hipchat.Message, error)
}

func NewClientMock() *ClientMock {

	cm := &ClientMock{}
	cm.sendMessage = func(string, string) error { return nil }
	cm.sendNotification = func(string, *hipchat.NotificationRequest) error { return nil }
	cm.getMessages = func(string, *hipchat.LatestHistoryOptions) ([]hipchat.Message, error) {
		return []hipchat.Message{}, nil
	}
	return cm
}

func (cm *ClientMock) SendMessage(roomID, message string) error {

	cm.SendMessageCall = SendMessageCall{roomId: roomID, message: message}
	return cm.sendMessage(roomID, message)
}

func (cm *ClientMock) SendNotification(roomID string, notification *hipchat.NotificationRequest) error {

	cm.SendNotificationCall = SendNotificationCall{roomId: roomID, notification: notification}
	return cm.sendNotification(roomID, notification)
}

func (cm *ClientMock) GetMessages(roomID string, options *hipchat.LatestHistoryOptions) ([]hipchat.Message, error) {

	cm.GetMessagesCall = GetMessagesCall{roomID: roomID, options: options}
	return cm.getMessages(roomID, options)
}
