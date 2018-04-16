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
	"github.com/stretchr/testify/assert"
	hc "github.com/tbruyelle/hipchat-go/hipchat"
	"github.com/HotelsDotCom/go-logger"
	"testing"
)

func TestToUserMessageUser(t *testing.T) {

	messageUser := map[string]interface{}{}
	messageUser["id"] = float64(81)
	messageUser["name"] = "John Rambo"
	messageUser["mention_name"] = "Rambo"

	user := ToUser(messageUser)

	assert.Equal(t, 81, user.Id)
	assert.Equal(t, "John Rambo", user.Name)
	assert.Equal(t, "Rambo", user.MentionName)
}

func TestToUserMessageUserInvalidFields(t *testing.T) {

	messageUser := map[string]interface{}{}
	messageUser["id"] = "81"
	messageUser["name"] = float64(99)
	messageUser["mention_name"] = 5

	user := ToUser(messageUser)

	assert.Equal(t, 0, user.Id)
	assert.Equal(t, "", user.Name)
	assert.Equal(t, "", user.MentionName)
}

func TestToUserHipchatUser(t *testing.T) {

	user := ToUser(hc.User{ID: 6, Name: "Karl Jr", MentionName: "Budha"})

	assert.Equal(t, 6, user.Id)
	assert.Equal(t, "Karl Jr", user.Name)
	assert.Equal(t, "Budha", user.MentionName)
}

func TestToUserString(t *testing.T) {

	user := ToUser("Carl W Jr Foox")

	assert.Equal(t, 0, user.Id)
	assert.Equal(t, "Carl W Jr Foox", user.Name)
	assert.Equal(t, "", user.MentionName)
}

func TestToUserInvalidType(t *testing.T) {

	l := NewMockLogger()
	defer l.rollback()

	user := ToUser(struct{ SomeField string }{"hello"})

	assert.Equal(t, 0, user.Id)
	assert.Equal(t, "", user.Name)
	assert.Equal(t, "", user.MentionName)
	assert.Contains(t, l.errorFMsg, "cannot convert ")
}

func TestToUserNil(t *testing.T) {

	user := ToUser(nil)

	assert.Equal(t, 0, user.Id)
	assert.Equal(t, "", user.Name)
	assert.Equal(t, "", user.MentionName)
}

type MockLogger struct {
	prevLogger func(string, ...interface{})
	errorFMsg  string
}

func NewMockLogger() *MockLogger {

	l := &MockLogger{}
	l.prevLogger = logger.Errorf
	logger.Errorf = l.Errorf
	return l
}

func (l *MockLogger) Errorf(format string, a ...interface{}) {
	l.errorFMsg = fmt.Sprintf(format, a...)
}

func (l *MockLogger) rollback() {
	logger.Errorf = l.prevLogger
}
