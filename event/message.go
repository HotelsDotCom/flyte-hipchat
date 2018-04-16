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

package event

import (
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/flyte-hipchat/hipchat"
	"github.com/HotelsDotCom/go-logger"
)

func HandleReceivedMessages(pack flyte.Pack, messages chan hipchat.Message) {

	go func() {
		for message := range messages {
			e := flyte.Event{
				EventDef: flyte.EventDef{Name: "ReceivedMessage"},
				Payload:  message,
			}
			logger.Infof("received message=%q in room=%s from=%q", message.Message, message.RoomId, message.From.Name)
			if err := pack.SendEvent(e); err != nil {
				logger.Errorf("error sending received message event: %v", err)
			}
		}
	}()
}
