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
	hc "github.com/tbruyelle/hipchat-go/hipchat"
	"log"
	"github.com/HotelsDotCom/flyte-hipchat/client"
	"time"
)

type Room struct {
	roomId        string
	client        client.HipchatClient
	messages      chan Message
	leave         chan bool
	lastMessageId string
}

func NewRoom(roomId string, client client.HipchatClient, messages chan Message) *Room {

	room := &Room{roomId: roomId, client: client, messages: messages, leave: make(chan bool)}
	room.monitor()
	return room
}

func (r *Room) Leave() {
	r.leave <- true
}

func (r *Room) monitor() {

	go func() {
		for {
			select {
			case <-r.leave:
				return
			default:
				r.handleIncomingMessages()
			}
		}
	}()
}

func (r *Room) handleIncomingMessages() {

	messages, err := r.getLatestMessages()
	if err != nil {
		log.Printf("cannot get room %s history: %v", r.roomId, err)
	}

	if len(messages) == 0 {
		time.Sleep(2 * time.Second)
		return
	}

	for _, message := range messages {
		r.lastMessageId = message.Id
		r.messages <- message
	}
}

func (r *Room) getLatestMessages() ([]Message, error) {

	if r.lastMessageId == "" {
		return r.getLatestMessage()
	}
	return r.getHistory()
}

func (r *Room) getLatestMessage() ([]Message, error) {

	messages, err := r.client.GetMessages(r.roomId, hipChatHistoryOptions("", 1))
	if err != nil {
		return []Message{}, err
	}
	return ToMessages(r.roomId, messages), nil
}

func (r *Room) getHistory() ([]Message, error) {

	messages, err := r.client.GetMessages(r.roomId, hipChatHistoryOptions(r.lastMessageId, 100))
	if err != nil {
		return []Message{}, err
	}

	if len(messages) > 1 {
		// response includes 'last message' at index 0
		return ToMessages(r.roomId, messages[1:]), nil
	}
	return []Message{}, nil
}

func hipChatHistoryOptions(fromId string, limit int) *hc.LatestHistoryOptions {

	return &hc.LatestHistoryOptions{
		MaxResults: limit,
		Timezone:   "GMT",
		NotBefore:  fromId,
	}
}
