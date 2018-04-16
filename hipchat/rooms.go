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
	"encoding/json"
	"io/ioutil"
	"github.com/HotelsDotCom/flyte-hipchat/client"
	"github.com/HotelsDotCom/go-logger"
	"sync"
)

type Rooms struct {
	sync.RWMutex
	backupPath string
	client     client.HipchatClient
	messages   chan Message
	rooms      map[string]*Room
}

func NewRooms(backupPath string, client client.HipchatClient, messages chan Message) (*Rooms, error) {

	r := &Rooms{
		backupPath: backupPath,
		client:     client,
		messages:   messages,
		rooms:      make(map[string]*Room),
	}

	if err := r.loadRooms(); err != nil {
		return r, err
	}
	return r, nil
}

func (r *Rooms) Add(roomId string) bool {

	r.Lock()
	defer r.Unlock()

	if _, ok := r.rooms[roomId]; !ok {
		r.rooms[roomId] = NewRoom(roomId, r.client, r.messages)
		r.save()
		return true
	}
	return false
}

func (r *Rooms) Remove(roomId string) {

	r.Lock()
	defer r.Unlock()

	if room, ok := r.rooms[roomId]; ok {
		go room.Leave()
		delete(r.rooms, roomId)
		r.save()
	}
}

func (r *Rooms) Get(roomId string) *Room {

	r.RLock()
	defer r.RUnlock()
	return r.rooms[roomId]
}

func (r *Rooms) ListIds() []string {

	r.RLock()
	defer r.RUnlock()

	rooms := []string{}
	for _, room := range r.rooms {
		rooms = append(rooms, room.roomId)
	}
	return rooms
}

func (r *Rooms) loadRooms() error {

	b, err := ioutil.ReadFile(r.backupPath)
	if err != nil {
		return err
	}

	roomIds := []string{}
	if err = json.Unmarshal(b, &roomIds); err != nil {
		return err
	}

	for _, id := range roomIds {
		r.rooms[id] = NewRoom(id, r.client, r.messages)
	}
	return nil
}

func (r *Rooms) save() {

	rooms := []string{}
	for k := range r.rooms {
		rooms = append(rooms, k)
	}

	b, err := json.Marshal(rooms)
	if err != nil {
		logger.Errorf("saving rooms, cannot marshal: %v", err)
		return
	}

	if err := ioutil.WriteFile(r.backupPath, b, 0644); err != nil {
		logger.Errorf("saving rooms, cannot write to file=%q: %v", r.backupPath, err)
	}
}
