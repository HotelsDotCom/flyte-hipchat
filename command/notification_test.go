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

package command

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/flyte-hipchat/hipchat"
	"testing"
)

func TestSendNotification(t *testing.T) {

	input := []byte(`{"roomId": "room id", "message": "test message", "from": "sender"}`)
	command := SendNotificationCommand(NewSendNotificationMock())
	event := command.Handler(input)

	expected := newNotificationEvent(SendNotificationOutput{
		SendNotificationInput: SendNotificationInput{RoomId: "room id", Message: "test message", From: "sender"},
	})
	assert.Equal(t, expected, event)
}

func TestSendNotificationCopiesInputToOutput(t *testing.T) {

	input := []byte(`{"roomId": "123", "message": "message", "from": "sender"}`)
	command := SendNotificationCommand(NewSendNotificationMock())
	event := command.Handler(input)

	expected := flyte.Event{
		EventDef: flyte.EventDef{Name: "NotificationSent"},
		Payload: SendNotificationOutput{
			SendNotificationInput: SendNotificationInput{RoomId: "123", Message: "message", From: "sender"},
		},
	}
	assert.Equal(t, expected, event)
}

func TestSendNotificationOptionalFields(t *testing.T) {

	input := []byte(`{"roomId": "456", "message": "message", "messageFormat": "html", "notify": true, "color": "blue", "from": "sender"}`)
	command := SendNotificationCommand(NewSendNotificationMock())
	event := command.Handler(input)

	expected := newNotificationEvent(SendNotificationOutput{
		SendNotificationInput: SendNotificationInput{
			RoomId:        "456",
			Message:       "message",
			MessageFormat: "html",
			Notify:        true,
			Color:         "blue",
			From:          "sender"}})
	assert.Equal(t, expected, event)
}

func TestSendNotificationMissingRequiredFields(t *testing.T) {

	cases := []struct {
		input         []byte
		expectedError string
	}{
		{[]byte(`{}`), "missing fields: [room id, message, from]"},
		{[]byte(`{"roomId": "123"}`), "missing fields: [message, from]"},
		{[]byte(`{"roomId": "123", "message": "the message"}`), "missing fields: [from]"},
		{[]byte(`{"message": "the message"}`), "missing fields: [room id, from]"},
		{[]byte(`{"message": "the message", "from": "from"}`), "missing fields: [room id]"},
		{[]byte(`{"from": "from"}`), "missing fields: [room id, message]"},
		{[]byte(`{"roomId": "xyz", "from": "from"}`), "missing fields: [message]"},
	}

	hc := NewSendNotificationMock()
	for _, c := range cases {

		command := SendNotificationCommand(hc)
		event := command.Handler(c.input)

		output := event.Payload.(SendNotificationErrorOutput)
		assert.Equal(t, c.expectedError, output.Error)
	}
}

func TestSendNotificationInvalidInput(t *testing.T) {

	input := []byte(`invalid message`)
	command := SendNotificationCommand(NewSendNotificationMock())
	event := command.Handler(input)

	output := event.Payload.(string)
	assert.Contains(t, output, "input is not valid:")
}

func TestSendNotificationMissingRoomId(t *testing.T) {

	input := []byte(`{"message": "the message", "from": "sender"}`)
	command := SendNotificationCommand(NewSendNotificationMock())
	event := command.Handler(input)

	output := SendNotificationOutput{
		SendNotificationInput: SendNotificationInput{
			Message: "the message",
			From:    "sender",
		},
	}
	expected := newNotificationFailedEvent(output, "missing fields: [room id]")
	assert.Equal(t, expected, event)
}

func TestSendNotificationMissingRoomAndMessage(t *testing.T) {

	input := []byte(`{"from": "sender"}`)
	command := SendNotificationCommand(NewSendNotificationMock())
	event := command.Handler(input)

	output := SendNotificationOutput{
		SendNotificationInput: SendNotificationInput{
			From: "sender",
		},
	}

	expected := newNotificationFailedEvent(output, "missing fields: [room id, message]")
	assert.Equal(t, expected, event)
}

func TestSendNotificationToHipchat(t *testing.T) {

	hc := NewSendNotificationMock()

	input := []byte(`{"roomId": "room id", "message": "test message", "from": "sender", "color": "pink"}`)
	command := SendNotificationCommand(hc)
	command.Handler(input)

	expected := hipchat.Notification{Message: "test message", From: "sender", Color: "pink"}
	assert.Equal(t, "room id", hc.RoomId)
	assert.Equal(t, expected, hc.Notification)
}

func TestSendNotificationToHipchatFailed(t *testing.T) {

	hc := NewSendNotificationMock()
	hc.sendNotification = func(string, hipchat.Notification) error { return errors.New("test error") }

	input := []byte(`{"roomId": "room id", "message": "test message", "from": "sender"}`)
	command := SendNotificationCommand(hc)
	event := command.Handler(input)

	output := SendNotificationOutput{
		SendNotificationInput: SendNotificationInput{
			RoomId:  "room id",
			Message: "test message",
			From:    "sender",
		},
	}

	expected := newNotificationFailedEvent(output, "error sending notification: test error")
	assert.Equal(t, expected, event)
}

func TestMarshalOutputEvent(t *testing.T) {

	input := []byte(`{"roomId": "room id", "message": "message", "messageFormat": "message format", "notify": true, "color": "red", "from": "Carl"}`)

	command := SendNotificationCommand(NewSendNotificationMock())
	event := command.Handler(input)

	payloadJson, _ := json.Marshal(event.Payload)
	expectedPayload := `{"roomId":"room id","message":"message","messageFormat":"message format","notify":true,"color":"red","from":"Carl"}`
	assert.Equal(t, expectedPayload, string(payloadJson))
}

func TestNotificationCommand(t *testing.T) {

	command := SendNotificationCommand(NewSendNotificationMock())

	assert.Equal(t, "SendNotification", command.Name)
	assert.Equal(t, 2, len(command.OutputEvents))
	assert.Equal(t, "NotificationSent", command.OutputEvents[0].Name)
	assert.Equal(t, "SendNotificationFailed", command.OutputEvents[1].Name)
}

type SendNotificationMock struct {
	RoomId           string
	Notification     hipchat.Notification
	sendNotification func(roomId string, notification hipchat.Notification) error
}

func NewSendNotificationMock() *SendNotificationMock {

	hc := &SendNotificationMock{}
	hc.sendNotification = func(roomId string, notification hipchat.Notification) error {
		hc.RoomId = roomId
		hc.Notification = notification
		return nil
	}
	return hc
}

func (hc *SendNotificationMock) SendNotification(roomId string, notification hipchat.Notification) error {
	return hc.sendNotification(roomId, notification)
}
