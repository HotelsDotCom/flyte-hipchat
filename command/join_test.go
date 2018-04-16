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

func TestGetJoinedEvent(t *testing.T) {

	payload := JoinRoomOutput{JoinRoomInput: JoinRoomInput{RoomId: "123"}}
	expected := flyte.Event{
		EventDef: flyte.EventDef{Name: "RoomJoined"},
		Payload:  payload,
	}

	actual := newJoinedEvent("123")

	assert.Equal(t, expected, actual)
}

func TestJoin(t *testing.T) {

	command := JoinCommand(NewHipchatRoomJoinerMock())
	expected := newJoinedEvent("123")

	event := command.Handler([]byte(`{"roomId": "123"}`))

	assert.Equal(t, expected, event)
}

func TestJoinMissingRoomIdField(t *testing.T) {

	command := JoinCommand(NewHipchatRoomJoinerMock())
	expected := newJoinedFailedEvent("", "missing room id field")

	event := command.Handler([]byte(`{}`))

	assert.Equal(t, expected, event)
}

func TestJoinInvalidInput(t *testing.T) {

	command := JoinCommand(NewHipchatRoomJoinerMock())

	event := command.Handler([]byte(`invalid input`))
	error := event.Payload.(string)

	assert.Contains(t, error, "input is not valid: ")
}

func TestJoinToHipchat(t *testing.T) {

	hc := NewHipchatRoomJoinerMock()

	command := JoinCommand(hc)
	command.Handler([]byte(`{"roomId": "456"}`))

	assert.Equal(t, "456", hc.CalledRoomId)

}

func TestJoinToHipchatFailed(t *testing.T) {

	hc := &HipchatRoomJoinerMock{}
	hc.joinRoom = func(string) error { return errors.New("test error") }
	expected := newJoinedFailedEvent("456", "cannot join room: test error")

	command := JoinCommand(hc)
	event := command.Handler([]byte(`{"roomId": "456"}`))

	assert.Equal(t, expected, event)
}

func TestJoinOutputEventMarshal(t *testing.T) {

	command := JoinCommand(NewHipchatRoomJoinerMock())

	event := command.Handler([]byte(`{"roomId": "xyz"}`))
	jsonPayload, _ := json.Marshal(event.Payload)

	assert.Equal(t, `{"roomId":"xyz"}`, string(jsonPayload))
}

func TestJoinCommand(t *testing.T) {

	command := JoinCommand(NewHipchatRoomJoinerMock())

	assert.Equal(t, "JoinRoom", command.Name)
	assert.Equal(t, 2, len(command.OutputEvents))
	assert.Equal(t, "RoomJoined", command.OutputEvents[0].Name)
	assert.Equal(t, "JoinRoomFailed", command.OutputEvents[1].Name)
}

type HipchatRoomJoinerMock struct {
	CalledRoomId string
	joinRoom     func(string) error
}

func NewHipchatRoomJoinerMock() *HipchatRoomJoinerMock {

	hc := &HipchatRoomJoinerMock{}
	hc.joinRoom = func(roomId string) error {
		hc.CalledRoomId = roomId
		return nil
	}
	return hc
}

func (hc *HipchatRoomJoinerMock) JoinRoom(roomId string) error {
	return hc.joinRoom(roomId)
}
