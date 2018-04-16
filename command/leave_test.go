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

func TestLeave(t *testing.T) {

	command := LeaveCommand(NewHipchatRoomLeaverMock())
	expected := flyte.Event{
		EventDef: flyte.EventDef{Name: "RoomLeft"},
		Payload:  LeaveRoomOutput{LeaveRoomInput: LeaveRoomInput{RoomId: "abc"}},
	}

	event := command.Handler([]byte(`{"roomId": "abc"}`))

	assert.Equal(t, expected, event)
}

func TestLeaveMissingRoomId(t *testing.T) {

	command := LeaveCommand(NewHipchatRoomLeaverMock())
	expected := newLeaveFailedEvent("", "missing room id field")

	event := command.Handler([]byte(`{}`))

	assert.Equal(t, expected, event)
}

func TestLeaveInvalidInput(t *testing.T) {

	command := LeaveCommand(NewHipchatRoomLeaverMock())

	event := command.Handler([]byte(`invalid input`))
	e := event.Payload.(string)

	assert.Contains(t, e, "input is not valid: ")
}

func TestSendLeaveToHipchat(t *testing.T) {

	hc := NewHipchatRoomLeaverMock()

	command := LeaveCommand(hc)
	command.Handler([]byte(`{"roomId": "987"}`))

	assert.Equal(t, "987", hc.CalledRoomId)
}

func TestSendLeaveToHipchatFailed(t *testing.T) {

	hc := NewHipchatRoomLeaverMock()
	hc.leaveRoom = func(string) error { return errors.New("test error") }

	expected := newLeaveFailedEvent("the room id", "cannot leave room: test error")
	command := LeaveCommand(hc)
	event := command.Handler([]byte(`{"roomId": "the room id"}`))

	assert.Equal(t, expected, event)
}

func TestLeaveOutputEventMarshal(t *testing.T) {

	command := LeaveCommand(NewHipchatRoomLeaverMock())

	event := command.Handler([]byte(`{"roomId": "xyz"}`))
	jsonPayload, _ := json.Marshal(event.Payload)

	assert.Equal(t, `{"roomId":"xyz"}`, string(jsonPayload))
}

func TestLeaveCommand(t *testing.T) {

	command := LeaveCommand(NewHipchatRoomLeaverMock())

	assert.Equal(t, "LeaveRoom", command.Name)
	assert.Equal(t, 2, len(command.OutputEvents))
	assert.Equal(t, "RoomLeft", command.OutputEvents[0].Name)
	assert.Equal(t, "LeaveRoomFailed", command.OutputEvents[1].Name)
}

type HipchatRoomLeaverMock struct {
	CalledRoomId string
	leaveRoom    func(string) error
}

func NewHipchatRoomLeaverMock() *HipchatRoomLeaverMock {

	hc := &HipchatRoomLeaverMock{}
	hc.leaveRoom = func(roomId string) error {
		hc.CalledRoomId = roomId
		return nil
	}
	return hc
}

func (hc *HipchatRoomLeaverMock) LeaveRoom(roomId string) error {
	return hc.leaveRoom(roomId)
}
