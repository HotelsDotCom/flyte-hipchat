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
	"github.com/HotelsDotCom/flyte-hipchat/hipchat"
	"strings"
)

type SendNotificationInput struct {
	RoomId        string `json:"roomId"`
	Message       string `json:"message"`
	MessageFormat string `json:"messageFormat"`
	Notify        bool   `json:"notify"`
	Color         string `json:"color"`
	From          string `json:"from"`
}

type SendNotificationOutput struct {
	SendNotificationInput
}

type SendNotificationErrorOutput struct {
	SendNotificationOutput
	Error string `json:"error"`
}

type HipchatNotificationSender interface {
	SendNotification(roomId string, notification hipchat.Notification) error
}

func SendNotificationCommand(hc HipchatNotificationSender) flyte.Command {

	return flyte.Command{
		Name:         "SendNotification",
		OutputEvents: []flyte.EventDef{{Name: "NotificationSent"}, {Name: "SendNotificationFailed"}},
		Handler:      sendNotificationHandler(hc),
	}
}

func sendNotificationHandler(hc HipchatNotificationSender) flyte.CommandHandler {

	return func(rawInput json.RawMessage) flyte.Event {

		input := SendNotificationInput{}
		if err := json.Unmarshal(rawInput, &input); err != nil {
			return flyte.NewFatalEvent(fmt.Sprintf("input is not valid: %v", err))
		}

		output := SendNotificationOutput{
			SendNotificationInput: input,
		}

		if err := validateNotification(input); err != nil {
			return newNotificationFailedEvent(output, err.Error())
		}

		if err := hc.SendNotification(input.RoomId, toClientNotification(input)); err != nil {
			return newNotificationFailedEvent(output, fmt.Sprintf("error sending notification: %v", err))
		}
		return newNotificationEvent(output)
	}
}

func toClientNotification(input SendNotificationInput) hipchat.Notification {

	return hipchat.Notification{
		Message:       input.Message,
		MessageFormat: input.MessageFormat,
		Notify:        input.Notify,
		Color:         input.Color,
		From:          input.From,
	}
}

func validateNotification(input SendNotificationInput) error {

	fields := []string{}
	if input.RoomId == "" {
		fields = append(fields, "room id")
	}
	if input.Message == "" {
		fields = append(fields, "message")
	}
	if input.From == "" {
		fields = append(fields, "from")
	}

	if len(fields) != 0 {
		return fmt.Errorf("missing fields: [%s]", strings.Join(fields, ", "))
	}
	return nil
}

func newNotificationEvent(output SendNotificationOutput) flyte.Event {

	return flyte.Event{
		EventDef: flyte.EventDef{Name: "NotificationSent"},
		Payload:  output,
	}
}

func newNotificationFailedEvent(output SendNotificationOutput, err string) flyte.Event {

	return flyte.Event{
		EventDef: flyte.EventDef{Name: "SendNotificationFailed"},
		Payload:  SendNotificationErrorOutput{SendNotificationOutput: output, Error: err},
	}
}
