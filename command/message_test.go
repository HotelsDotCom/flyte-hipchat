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
	"testing"
)

func TestGetEventCopiesOutputToFlyteEventAndSetsEventDefName(t *testing.T) {

	event := newMessageSentEvent("room id", "the message")

	expected := flyte.Event{
		EventDef: flyte.EventDef{Name: "MessageSent"},
		Payload:  SendMessageOutput{SendMessageInput: SendMessageInput{RoomId: "room id", Message: "the message"}},
	}
	assert.Equal(t, expected, event)
}

func TestGetErrorEventCopiesOutputToFlyteEventAndSetsEventDefName(t *testing.T) {

	event := newSendMessageFailedEvent("room id", "the message", "the error")

	output := SendMessageOutput{SendMessageInput: SendMessageInput{RoomId: "room id", Message: "the message"}}
	expected := flyte.Event{
		EventDef: flyte.EventDef{Name: "SendMessageFailed"},
		Payload:  SendMessageErrorOutput{SendMessageOutput: output, Error: "the error"},
	}
	assert.Equal(t, expected, event)
}

func TestSendMessage(t *testing.T) {

	hcInput := struct {
		sentRoomId  string
		sentMessage string
	}{}

	hc := getHCMock()
	hc.sendMessage = func(roomId string, message string) error {
		hcInput.sentRoomId = roomId
		hcInput.sentMessage = message
		return nil
	}

	command := SendMessageCommand(hc)
	command.Handler([]byte(`{"roomId": "123", "message": "Hello"}`))

	assert.Equal(t, "123", hcInput.sentRoomId)
	assert.Equal(t, "Hello", hcInput.sentMessage)
}

func TestSendMessageFailure(t *testing.T) {

	hc := getHCMock()
	hc.sendMessage = func(roomId string, message string) error {
		return errors.New("the error")
	}

	command := SendMessageCommand(hc)
	event := command.Handler([]byte(`{"roomId": "123", "message": "Hello"}`))

	expectedEvent := newSendMessageFailedEvent("123", "Hello", "error sending message: the error")

	assert.Equal(t, "SendMessageFailed", event.EventDef.Name)
	assert.Equal(t, "SendMessageFailed", expectedEvent.EventDef.Name)
	assert.Equal(t, expectedEvent, event)
}

func TestSendMessageSuccessOutputAlsoContainsInput(t *testing.T) {

	command := SendMessageCommand(getHCMock())
	event := command.Handler([]byte(`{"roomId": "123", "message": "Hello"}`))

	expected := newMessageSentEvent("123", "Hello")
	assert.Equal(t, expected, event)
}

func TestSendInvalidMessage(t *testing.T) {

	command := SendMessageCommand(getHCMock())
	event := command.Handler([]byte(`invalid message`))

	output := event.Payload.(string)
	assert.Contains(t, output, "input is not valid:")
}

func TestSendNilMessage(t *testing.T) {

	command := SendMessageCommand(getHCMock())
	event := command.Handler(nil)

	output := event.Payload.(string)
	assert.Contains(t, output, "input is not valid:")
}

func TestSendMessageWithMissingFields(t *testing.T) {

	cases := []struct {
		input         SendMessageInput
		expectedError string
	}{
		{input: SendMessageInput{}, expectedError: "missing fields: [room id, message]"},
		{input: SendMessageInput{Message: "the message"}, expectedError: "missing fields: [room id]"},
		{input: SendMessageInput{RoomId: "the room id"}, expectedError: "missing fields: [message]"},
		{input: SendMessageInput{RoomId: "the room id", Message: "the message"}, expectedError: ""},
	}

	command := SendMessageCommand(getHCMock())
	for _, c := range cases {
		b, _ := json.Marshal(c.input)
		event := command.Handler(b)
		expected := newMessageSentEvent(c.input.RoomId, c.input.Message)
		if c.expectedError != "" {
			expected = newSendMessageFailedEvent(c.input.RoomId, c.input.Message, c.expectedError)
		}
		assert.Equal(t, expected, event)
	}
}

func TestMarshalSendMessageOutput(t *testing.T) {

	command := SendMessageCommand(getHCMock())
	event := command.Handler([]byte(`{"roomId": "123", "message": "Hello"}`))
	payloadJson, _ := json.Marshal(event.Payload)

	assert.Equal(t, `{"roomId":"123","message":"Hello"}`, string(payloadJson))
}

func TestMessageCommand(t *testing.T) {

	command := SendMessageCommand(getHCMock())

	assert.Equal(t, "SendMessage", command.Name)
	assert.Equal(t, 2, len(command.OutputEvents))
	assert.Equal(t, "MessageSent", command.OutputEvents[0].Name)
	assert.Equal(t, "SendMessageFailed", command.OutputEvents[1].Name)
}

func TestMarshalOutputContainingError(t *testing.T) {

	hc := getHCMock()
	hc.sendMessage = func(roomId string, message string) error {
		return errors.New("the error")
	}

	command := SendMessageCommand(hc)
	event := command.Handler([]byte(`{"roomId": "123", "message": "Hello"}`))
	payloadJson, _ := json.Marshal(event.Payload) // event is handled by client, only payload is ours

	expectedJson := `{"roomId":"123","message":"Hello","error":"error sending message: the error"}`
	assert.Equal(t, expectedJson, string(payloadJson))
}

func TestCommandDefinition(t *testing.T) {

	command := SendMessageCommand(nil)

	assert.Equal(t, "SendMessage", command.Name)
	assert.Len(t, command.OutputEvents, 2)
	assert.Equal(t, "MessageSent", command.OutputEvents[0].Name)
	assert.Equal(t, "SendMessageFailed", command.OutputEvents[1].Name)
	assert.Nil(t, command.OutputEvents[0].HelpURL)
	assert.Nil(t, command.HelpURL)
}

type hipchatMessageSenderMock struct {
	sendMessage func(string, string) error
}

func getHCMock() hipchatMessageSenderMock {
	return hipchatMessageSenderMock{sendMessage: func(string, string) error { return nil }}
}

func (hm hipchatMessageSenderMock) SendMessage(roomId string, message string) error {
	return hm.sendMessage(roomId, message)
}
