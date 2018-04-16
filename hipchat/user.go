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
	hc "github.com/tbruyelle/hipchat-go/hipchat"
	"github.com/HotelsDotCom/go-logger"
)

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	MentionName string `json:"mentionName"`
}

func ToUser(user interface{}) User {

	if user == nil {
		return User{}
	}

	// hipchat user
	if u, ok := user.(hc.User); ok {
		return fromHipchatUser(u)
	}

	// notification - only username
	if name, ok := user.(string); ok {
		return User{Name: name}
	}

	// message
	if from, ok := user.(map[string]interface{}); ok {
		return fromMessageUser(from)
	}

	logger.Errorf("cannot convert 'hipchat User/From' %+v", user)
	return User{}
}

func fromHipchatUser(user hc.User) User {

	return User{
		Id:          user.ID,
		Name:        user.Name,
		MentionName: user.MentionName,
	}
}

func fromMessageUser(user map[string]interface{}) User {

	u := User{}
	if id, ok := user["id"]; ok {
		// hipchat client inconsistencies, user id in hipchat.User struct is int, in here it is float64
		if v, ok := id.(float64); ok {
			u.Id = int(v)
		}
	}
	if mentionName, ok := user["mention_name"]; ok {
		if v, ok := mentionName.(string); ok {
			u.MentionName = v
		}
	}
	if name, ok := user["name"]; ok {
		if v, ok := name.(string); ok {
			u.Name = v

		}
	}
	return u
}
