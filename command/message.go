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
	"strings"
)

type SendMessageInput struct {
	RoomId  string `json:"roomId"`
	Message string `json:"message"`
}

type SendMessageOutput struct {
	SendMessageInput
}

type SendMessageErrorOutput struct {
	SendMessageOutput
	Error string `json:"error"`
}

type HipchatMessageSender interface {
	SendMessage(roomId string, message string) error
}

func SendMessageCommand(hc HipchatMessageSender) flyte.Command {

	return flyte.Command{
		Name:         "SendMessage",
		OutputEvents: []flyte.EventDef{{Name: "MessageSent"}, {Name: "SendMessageFailed"}},
		Handler:      sendMessageHandler(hc),
	}
}

func sendMessageHandler(hc HipchatMessageSender) flyte.CommandHandler {

	return func(rawInput json.RawMessage) flyte.Event {

		input := SendMessageInput{}
		if err := json.Unmarshal(rawInput, &input); err != nil {
			return flyte.NewFatalEvent(fmt.Sprintf("input is not valid: %v", err))
		}

		if err := validateMessage(input); err != nil {
			return newSendMessageFailedEvent(input.RoomId, input.Message, err.Error())
		}

		if err := hc.SendMessage(input.RoomId, input.Message); err != nil {
			return newSendMessageFailedEvent(input.RoomId, input.Message, fmt.Sprintf("error sending message: %v", err))
		}
		return newMessageSentEvent(input.RoomId, input.Message)
	}
}

func validateMessage(input SendMessageInput) error {

	errors := []string{}
	if input.RoomId == "" {
		errors = append(errors, "room id")
	}
	if input.Message == "" {
		errors = append(errors, "message")
	}
	if len(errors) != 0 {
		return fmt.Errorf("missing fields: [%s]", strings.Join(errors, ", "))
	}
	return nil
}

func newMessageSentEvent(roomId, message string) flyte.Event {

	return flyte.Event{
		EventDef: flyte.EventDef{Name: "MessageSent"},
		Payload:  SendMessageOutput{SendMessageInput: SendMessageInput{RoomId: roomId, Message: message}},
	}
}

func newSendMessageFailedEvent(roomId, message, err string) flyte.Event {

	output := SendMessageOutput{SendMessageInput: SendMessageInput{RoomId: roomId, Message: message}}
	return flyte.Event{
		EventDef: flyte.EventDef{Name: "SendMessageFailed"},
		Payload:  SendMessageErrorOutput{SendMessageOutput: output, Error: err},
	}
}
