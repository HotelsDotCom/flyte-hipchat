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
	"strings"
)

type Notification struct {
	Message       string
	MessageFormat string
	Notify        bool
	Color         string
	From          string
}

type Message struct {
	Id            string `json:"id"`
	RoomId        string `json:"roomId"`
	Date          string `json:"date"`
	From          User   `json:"from"`
	Mentions      []User `json:"mentions"`
	Message       string `json:"message"`
	MessageFormat string `json:"messageFormat"`
	Type          string `json:"type"`
}

func ToMessages(roomId string, hipchatMessages []hc.Message) []Message {

	messages := []Message{}
	for _, m := range hipchatMessages {
		messages = append(messages, toMessage(roomId, m))
	}
	return messages
}

func toMessage(roomId string, hipchatMessage hc.Message) Message {

	mentions := []User{}
	for _, hcUser := range hipchatMessage.Mentions {
		mentions = append(mentions, ToUser(hcUser))
	}

	return Message{
		Id:            hipchatMessage.ID,
		RoomId:        roomId,
		Date:          hipchatMessage.Date,
		From:          ToUser(hipchatMessage.From),
		Mentions:      mentions,
		Message:       hipchatMessage.Message,
		MessageFormat: hipchatMessage.MessageFormat,
		Type:          hipchatMessage.Type,
	}
}

func ToHipChatNotification(notification Notification) *hc.NotificationRequest {

	return &hc.NotificationRequest{
		Message:       notification.Message,
		MessageFormat: notification.MessageFormat,
		Notify:        notification.Notify,
		Color:         strToColor(notification.Color),
		From:          notification.From,
	}
}

func strToColor(color string) hc.Color {

	colors := map[string]hc.Color{
		"yellow": hc.ColorYellow,
		"green":  hc.ColorGreen,
		"red":    hc.ColorRed,
		"purple": hc.ColorPurple,
		"gray":   hc.ColorGray,
		"random": hc.ColorRandom,
	}

	if c, ok := colors[strings.ToLower(color)]; ok {
		return c
	}
	return hc.ColorYellow
}
