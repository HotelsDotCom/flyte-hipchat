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
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
)

func TestAddAndRemoveRoom(t *testing.T) {

	path := roomsPath()
	defer func() { os.Remove(path) }()

	rooms, _ := NewRooms(path, NewHipchatClientMock(), nil)
	assert.Equal(t, 0, len(rooms.ListIds()))

	rooms.Add("test room")
	assert.Equal(t, 1, len(rooms.ListIds()))
	assert.Equal(t, "test room", rooms.Get("test room").roomId)

	rooms.Remove("test room")
	assert.Equal(t, 0, len(rooms.ListIds()))
}

func TestAddRoomTwice(t *testing.T) {

	path := roomsPath()
	defer func() { os.Remove(path) }()

	rooms, _ := NewRooms(path, NewHipchatClientMock(), nil)
	assert.Equal(t, 0, len(rooms.ListIds()))

	ok := rooms.Add("test room")
	assert.True(t, ok)
	assert.Equal(t, 1, len(rooms.ListIds()))
	assert.Equal(t, "test room", rooms.Get("test room").roomId)

	ok = rooms.Add("test room")
	assert.False(t, ok)
	assert.Equal(t, 1, len(rooms.ListIds()))
	assert.Equal(t, "test room", rooms.Get("test room").roomId)
}

func TestConcurrentlyRemoveRooms(t *testing.T) {

	path := roomsPath()
	defer func() { os.Remove(path) }()

	rooms, _ := NewRooms(path, NewHipchatClientMock(), nil)
	for i := 0; i < 501; i++ {
		rooms.Add(strconv.Itoa(i))
	}

	var wg sync.WaitGroup
	wg.Add(500)
	for i := 0; i < 500; i++ {
		go func(i int) {
			rooms.Remove(strconv.Itoa(i))
			defer wg.Done()
		}(i)
	}

	wg.Wait()
	assert.Equal(t, 1, len(rooms.ListIds()))
	assert.Equal(t, "500", rooms.ListIds()[0])
	assert.Equal(t, "500", rooms.Get("500").roomId)
}

func TestConcurrentlyAddRooms(t *testing.T) {

	path := roomsPath()
	defer func() { os.Remove(path) }()

	rooms, _ := NewRooms(path, NewHipchatClientMock(), nil)

	var wg sync.WaitGroup
	wg.Add(500)
	for i := 0; i < 500; i++ {
		go func(i int) {
			defer wg.Done()
			rooms.Add(strconv.Itoa(i))
		}(i)
	}

	wg.Wait()
	assert.Equal(t, 500, len(rooms.ListIds()))
}

func TestRoomsPersistence(t *testing.T) {

	path := roomsPath()
	defer func() { os.Remove(path) }()

	rooms, _ := NewRooms(path, NewHipchatClientMock(), nil)

	var wg sync.WaitGroup
	wg.Add(500)
	for i := 0; i < 500; i++ {
		go func(i int) {
			defer wg.Done()
			rooms.Add(strconv.Itoa(i))
		}(i)
	}
	wg.Wait()

	rooms2, _ := NewRooms(path, NewHipchatClientMock(), nil)
	assert.Equal(t, 500, len(rooms2.ListIds()))
	assert.Equal(t, "483", rooms2.Get("483").roomId)
}

func roomsPath() string {

	path, _ := filepath.Abs(filepath.Dir(filepath.Join(os.Args[0], "test_rooms.json")))
	return path
}

type HipchatClientMock struct {
	sendMessage      func(roomID, message string) error
	sendNotification func(roomID string, notification *hipchat.NotificationRequest) error
	getMessages      func(roomID string, options *hipchat.LatestHistoryOptions) ([]hipchat.Message, error)
}

func NewHipchatClientMock() HipchatClientMock {

	hc := HipchatClientMock{}
	hc.sendMessage = func(string, string) error { return nil }
	hc.sendNotification = func(string, *hipchat.NotificationRequest) error { return nil }
	hc.getMessages = func(string, *hipchat.LatestHistoryOptions) ([]hipchat.Message, error) {
		return []hipchat.Message{}, nil
	}
	return hc
}

func (hc HipchatClientMock) SendMessage(roomID, message string) error {
	return hc.sendMessage(roomID, message)
}

func (hc HipchatClientMock) SendNotification(roomID string, notification *hipchat.NotificationRequest) error {
	return hc.sendNotification(roomID, notification)
}

func (hc HipchatClientMock) GetMessages(roomID string, options *hipchat.LatestHistoryOptions) ([]hipchat.Message, error) {
	return hc.getMessages(roomID, options)
}
