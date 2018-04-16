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

package hipchat

import (
	"fmt"
	"github.com/HotelsDotCom/flyte-hipchat/client"
	"github.com/HotelsDotCom/go-logger"
)

type Hipchat struct {
	client client.HipchatClient
	rooms  *Rooms
}

func NewHipchat(roomsBackupPath string, client client.HipchatClient, messages chan Message) (Hipchat, error) {

	hc := Hipchat{client: client}
	rooms, err := NewRooms(roomsBackupPath, client, messages)
	if err != nil {
		return hc, err
	}

	hc.rooms = rooms
	for _, id := range hc.rooms.ListIds() {
		if err := hc.SendNotification(id, startUpNotification); err != nil {
			logger.Errorf("room=%s cannot send startup notification: %v", id, err)
		}
	}

	return hc, nil
}

func (hc Hipchat) JoinedRoomIds() []string {
	return hc.rooms.ListIds()
}

func (hc Hipchat) BroadcastMessage(message string) error {

	errors := []string{}
	logger.Info("broadcasting message")
	for _, id := range hc.rooms.ListIds() {
		if err := hc.SendMessage(id, message); err != nil {
			errors = append(errors, fmt.Sprintf("room=%s: %v", id, err))
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("failed messages: %v", errors)
	}
	return nil
}

func (hc Hipchat) SendMessage(roomId string, message string) error {
	return hc.client.SendMessage(roomId, message)
}

func (hc Hipchat) SendNotification(roomId string, notification Notification) error {
	return hc.client.SendNotification(roomId, ToHipChatNotification(notification))
}

func (hc Hipchat) JoinRoom(roomId string) error {

	logger.Infof("joining room=%s", roomId)
	if !hc.rooms.Add(roomId) {
		logger.Infof("room=%s already joined", roomId)
		return nil
	}

	if r := hc.rooms.Get(roomId); r != nil {
		if err := hc.SendNotification(r.roomId, joinNotification); err != nil {
			return fmt.Errorf("cannot send notification to room=%s: %v", roomId, err)
		}
	}
	return nil
}

func (hc Hipchat) LeaveRoom(roomId string) error {

	var err error
	if r := hc.rooms.Get(roomId); r != nil {
		if e := hc.SendNotification(r.roomId, leaveNotification); e != nil {
			err = fmt.Errorf("cannot send notification to room=%s: %v", roomId, e)
		}
		logger.Infof("leaving room=%s", roomId)
		hc.rooms.Remove(roomId)
	}
	return err
}

func (hc Hipchat) Shutdown() {

	for _, id := range hc.rooms.ListIds() {

		room := hc.rooms.Get(id)
		room.Leave()
		if err := hc.SendNotification(room.roomId, shutDownNotification); err != nil {
			logger.Errorf("room=%s cannot send shutdown notification: %v", id, err)
		}
	}
	close(hc.rooms.messages)
}
