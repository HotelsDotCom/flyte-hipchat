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
	"fmt"
	"github.com/HotelsDotCom/flyte-client/flyte"
)

type LeaveRoomInput struct {
	RoomId string `json:"roomId"`
}

type LeaveRoomOutput struct {
	LeaveRoomInput
}

type LeaveRoomErrorOutput struct {
	LeaveRoomOutput
	Error string `json:"error"`
}

type HipchatRoomLeaver interface {
	LeaveRoom(roomId string) error
}

func LeaveCommand(hc HipchatRoomLeaver) flyte.Command {

	return flyte.Command{
		Name:         "LeaveRoom",
		OutputEvents: []flyte.EventDef{{Name: "RoomLeft"}, {Name: "LeaveRoomFailed"}},
		Handler:      leaveHandler(hc),
	}
}

func leaveHandler(hc HipchatRoomLeaver) flyte.CommandHandler {

	return func(rawInput json.RawMessage) flyte.Event {

		input := LeaveRoomInput{}
		if err := json.Unmarshal(rawInput, &input); err != nil {
			return flyte.NewFatalEvent(fmt.Sprintf("input is not valid: %v", err))
		}

		if input.RoomId == "" {
			return newLeaveFailedEvent(input.RoomId, "missing room id field")
		}

		if err := hc.LeaveRoom(input.RoomId); err != nil {
			return newLeaveFailedEvent(input.RoomId, fmt.Sprintf("cannot leave room: %v", err))
		}
		return newLeaveEvent(input.RoomId)
	}
}

func newLeaveEvent(roomId string) flyte.Event {

	return flyte.Event{
		EventDef: flyte.EventDef{Name: "RoomLeft"},
		Payload:  LeaveRoomOutput{LeaveRoomInput: LeaveRoomInput{RoomId: roomId}},
	}
}

func newLeaveFailedEvent(roomId, err string) flyte.Event {

	output := LeaveRoomOutput{LeaveRoomInput: LeaveRoomInput{RoomId: roomId}}
	return flyte.Event{
		EventDef: flyte.EventDef{Name: "LeaveRoomFailed"},
		Payload:  LeaveRoomErrorOutput{LeaveRoomOutput: output, Error: err},
	}
}
