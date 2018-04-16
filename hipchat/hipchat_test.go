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
	"github.com/stretchr/testify/assert"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/HotelsDotCom/flyte-hipchat/bkp"
	"github.com/HotelsDotCom/go-logger"
	"testing"
)

func TestNewHipchatNoBackedUpRooms(t *testing.T) {

	bkpPath := bkp.CreateBkpFile(createTestBkpDir(), "rooms.json")
	defer func() { os.Remove(bkpPath) }()

	notifiedRooms := []string{}
	client := NewClientMock()
	client.sendNotification = func(roomId string, notification *hipchat.NotificationRequest) error {
		notifiedRooms = append(notifiedRooms, roomId)
		return nil
	}

	hc, _ := NewHipchat(bkpPath, client, nil)

	assert.Equal(t, 0, len(hc.JoinedRoomIds()))
	assert.Equal(t, 0, len(notifiedRooms))
}

func TestNewHipchatBackedUpRooms(t *testing.T) {

	bkpPath := bkp.CreateBkpFile(createTestBkpDir(), "rooms.json")
	defer func() { os.Remove(bkpPath) }()
	ioutil.WriteFile(bkpPath, []byte(`["123", "456"]`), 0644)

	notifiedRooms := []string{}
	client := NewClientMock()
	client.sendNotification = func(roomId string, notification *hipchat.NotificationRequest) error {
		notifiedRooms = append(notifiedRooms, roomId)
		return nil
	}

	hc, _ := NewHipchat(bkpPath, client, nil)

	joinedRooms := hc.JoinedRoomIds()
	assert.Equal(t, 2, len(joinedRooms))
	assert.Contains(t, joinedRooms, "123")
	assert.Contains(t, joinedRooms, "456")

	assert.Equal(t, 2, len(notifiedRooms))
	assert.Contains(t, notifiedRooms, "123")
	assert.Contains(t, notifiedRooms, "456")
}

func TestBroadcastMessage(t *testing.T) {

	bkpPath := bkp.CreateBkpFile(createTestBkpDir(), "rooms.json")
	defer func() { os.Remove(bkpPath) }()

	notifiedRooms := []string{}
	client := NewClientMock()
	client.sendMessage = func(roomId string, message string) error {
		notifiedRooms = append(notifiedRooms, roomId)
		return nil
	}

	hc, _ := NewHipchat(bkpPath, client, nil)

	assert.Equal(t, 0, len(notifiedRooms))
	hc.JoinRoom("1")
	hc.JoinRoom("2")

	hc.BroadcastMessage("test broadcast message")
	assert.Equal(t, 2, len(notifiedRooms))
	assert.Contains(t, notifiedRooms, "1")
	assert.Contains(t, notifiedRooms, "2")
}

func TestSendMessage(t *testing.T) {

	bkpPath := bkp.CreateBkpFile(createTestBkpDir(), "rooms.json")
	defer func() { os.Remove(bkpPath) }()

	notifiedRooms := []string{}
	client := NewClientMock()
	client.sendMessage = func(roomId string, message string) error {
		notifiedRooms = append(notifiedRooms, roomId)
		return nil
	}

	hc, _ := NewHipchat(bkpPath, client, nil)

	hc.SendMessage("the room id", "the message")
	assert.Equal(t, 1, len(notifiedRooms))
	assert.Equal(t, "the room id", notifiedRooms[0])
}

func TestSendNotification(t *testing.T) {

	bkpPath := bkp.CreateBkpFile(createTestBkpDir(), "rooms.json")
	defer func() { os.Remove(bkpPath) }()

	notifiedRooms := []string{}
	client := NewClientMock()
	client.sendNotification = func(roomId string, notification *hipchat.NotificationRequest) error {
		notifiedRooms = append(notifiedRooms, roomId)
		return nil
	}

	hc, _ := NewHipchat(bkpPath, client, nil)

	hc.SendNotification("the room id", joinNotification)
	assert.Equal(t, 1, len(notifiedRooms))
	assert.Equal(t, "the room id", notifiedRooms[0])
}

func TestLeaveRoom(t *testing.T) {

	bkpPath := bkp.CreateBkpFile(createTestBkpDir(), "rooms.json")
	defer func() { os.Remove(bkpPath) }()

	client := NewClientMock()

	hc, _ := NewHipchat(bkpPath, client, nil)

	assert.Equal(t, 0, len(hc.JoinedRoomIds()))

	hc.JoinRoom("1")
	hc.JoinRoom("2")
	assert.Equal(t, 2, len(hc.JoinedRoomIds()))

	hc.LeaveRoom("1")
	assert.Equal(t, 1, len(hc.JoinedRoomIds()))
	assert.Equal(t, "2", hc.JoinedRoomIds()[0])
}

func createTestBkpDir() string {

	dir := filepath.Join(os.TempDir(), "flyte-test-hipchat")
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Fatalf("cannot create bkp dir %s: %v", dir, err)
	}
	return dir
}
