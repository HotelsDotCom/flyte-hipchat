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

func TestBroadcastMessage(t *testing.T) {

	input := []byte(`{"message": "the message"}`)
	command := BroadcastCommand(NewHipchatBroadcasterMock())
	event := command.Handler(input)

	payload := BroadcastOutput{BroadcastInput: BroadcastInput{Message: "the message"}}
	expected := flyte.Event{EventDef: flyte.EventDef{Name: "BroadcastSent"}, Payload: payload}

	assert.Equal(t, expected, event)
}

func TestBroadcastMessageMissingMessage(t *testing.T) {

	input := []byte(`{}`)
	command := BroadcastCommand(NewHipchatBroadcasterMock())
	event := command.Handler(input)

	output := event.Payload.(BroadcastErrorOutput)
	assert.Equal(t, "missing message field", output.Error)
}

func TestBroadcastMessageInvalidInput(t *testing.T) {

	input := []byte(`invalid input`)
	command := BroadcastCommand(NewHipchatBroadcasterMock())
	event := command.Handler(input)

	// fatal event
	output := event.Payload.(string)
	assert.Contains(t, output, "input is not valid")
}

func TestBroadcastMessageToHipchat(t *testing.T) {

	hc := HipchatBroadcasterMock{}
	var receivedMessage string
	hc.broadcastMessage = func(message string) error { receivedMessage = message; return nil }

	input := []byte(`{"message": "the message"}`)
	command := BroadcastCommand(hc)
	command.Handler(input)

	assert.Equal(t, "the message", receivedMessage)
}

func TestBroadcastMessageToHipchatFailed(t *testing.T) {

	hc := HipchatBroadcasterMock{}
	hc.broadcastMessage = func(string) error { return errors.New("test error") }

	input := []byte(`{"message": "the message"}`)
	command := BroadcastCommand(hc)
	event := command.Handler(input)

	expected := newBroadcastFailedEvent("the message", "error broadcasting message: test error")
	assert.Equal(t, expected, event)
}

func TestMarshalSuccessBroadcastOutput(t *testing.T) {

	input := []byte(`{"message": "the message"}`)
	command := BroadcastCommand(NewHipchatBroadcasterMock())
	event := command.Handler(input)

	jsonEvent, _ := json.Marshal(event.Payload)
	assert.Equal(t, `{"message":"the message"}`, string(jsonEvent))
}

func TestMarshalErrorOutput(t *testing.T) {

	hc := HipchatBroadcasterMock{}
	hc.broadcastMessage = func(string) error { return errors.New("the error") }

	input := []byte(`{"message": "the message"}`)
	command := BroadcastCommand(hc)
	event := command.Handler(input)

	jsonPayload, _ := json.Marshal(event.Payload)
	assert.Equal(t, `{"message":"the message","error":"error broadcasting message: the error"}`, string(jsonPayload))
}

func TestBroadcastCommand(t *testing.T) {

	command := BroadcastCommand(NewHipchatBroadcasterMock())

	assert.Equal(t, "Broadcast", command.Name)
	assert.Equal(t, 2, len(command.OutputEvents))
	assert.Equal(t, "BroadcastSent", command.OutputEvents[0].Name)
	assert.Equal(t, "BroadcastFailed", command.OutputEvents[1].Name)
}

type HipchatBroadcasterMock struct {
	broadcastMessage func(string) error
}

func NewHipchatBroadcasterMock() HipchatBroadcasterMock {

	hc := HipchatBroadcasterMock{}
	hc.broadcastMessage = func(string) error { return nil }
	return hc
}

func (hc HipchatBroadcasterMock) BroadcastMessage(message string) error {
	return hc.broadcastMessage(message)
}
