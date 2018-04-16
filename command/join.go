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

type JoinRoomInput struct {
	RoomId string `json:"roomId"`
}

type JoinRoomOutput struct {
	JoinRoomInput
}

type JoinRoomErrorOutput struct {
	JoinRoomOutput
	Error string `json:"error"`
}

type HipchatRoomJoiner interface {
	JoinRoom(roomId string) error
}

func JoinCommand(hc HipchatRoomJoiner) flyte.Command {

	return flyte.Command{
		Name:         "JoinRoom",
		OutputEvents: []flyte.EventDef{{Name: "RoomJoined"}, {Name: "JoinRoomFailed"}},
		Handler:      joinHandler(hc),
	}
}

func joinHandler(hc HipchatRoomJoiner) flyte.CommandHandler {

	return func(rawInput json.RawMessage) flyte.Event {

		input := JoinRoomInput{}
		if err := json.Unmarshal(rawInput, &input); err != nil {
			return flyte.NewFatalEvent(fmt.Sprintf("input is not valid: %v", err))
		}

		if input.RoomId == "" {
			return newJoinedFailedEvent(input.RoomId, "missing room id field")
		}

		if err := hc.JoinRoom(input.RoomId); err != nil {
			return newJoinedFailedEvent(input.RoomId, fmt.Sprintf("cannot join room: %v", err))
		}
		return newJoinedEvent(input.RoomId)
	}
}

func newJoinedEvent(roomId string) flyte.Event {

	return flyte.Event{
		EventDef: flyte.EventDef{Name: "RoomJoined"},
		Payload:  JoinRoomOutput{JoinRoomInput: JoinRoomInput{RoomId: roomId}},
	}
}

func newJoinedFailedEvent(roomId, err string) flyte.Event {

	output := JoinRoomOutput{JoinRoomInput: JoinRoomInput{RoomId: roomId}}
	return flyte.Event{
		EventDef: flyte.EventDef{Name: "JoinRoomFailed"},
		Payload:  JoinRoomErrorOutput{JoinRoomOutput: output, Error: err},
	}
}
