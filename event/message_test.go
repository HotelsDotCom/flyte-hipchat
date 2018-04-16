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
	"github.com/stretchr/testify/assert"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/flyte-hipchat/hipchat"
	"testing"
)

func TestMessageReceived(t *testing.T) {

	p := NewPackMock()
	messages := make(chan hipchat.Message)
	HandleReceivedMessages(p, messages)

	messages <- hipchat.Message{Message: "the message"}
	receivedEvent := <-p.receivedEvents
	receivedEventDef := receivedEvent.EventDef
	receivedPayload := receivedEvent.Payload.(hipchat.Message)

	assert.Equal(t, "ReceivedMessage", receivedEventDef.Name)
	assert.Equal(t, "the message", receivedPayload.Message)
}

type PackMock struct {
	receivedEvents chan flyte.Event
	sendEvent      func(flyte.Event) error
}

func NewPackMock() *PackMock {

	p := &PackMock{}
	p.receivedEvents = make(chan flyte.Event)
	p.sendEvent = func(event flyte.Event) error {
		p.receivedEvents <- event
		return nil
	}
	return p
}

func (p *PackMock) Start() {
	return
}

func (p *PackMock) SendEvent(event flyte.Event) error {
	return p.sendEvent(event)
}
