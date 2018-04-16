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

package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/HotelsDotCom/go-logger"
	"testing"
)

func TestApiHost(t *testing.T) {

	os.Setenv("FLYTE_API", "http://test_api:8080")
	defer func() { os.Unsetenv("FLYTE_API") }()

	url := ApiHost()
	assert.Equal(t, "http://test_api:8080", url.String())
}

func TestApiHostNotSet(t *testing.T) {

	mockLogger := NewMockLogger()
	defer func() { mockLogger.rollback() }()

	ApiHost()
	assert.Equal(t, "env=FLYTE_API not set", mockLogger.fatalFMsg)
}

func TestApiHostInvalidUrl(t *testing.T) {

	os.Setenv("FLYTE_API", ":/invalid url")
	defer func() { os.Unsetenv("FLYTE_API") }()

	mockLogger := NewMockLogger()
	defer func() { mockLogger.rollback() }()

	ApiHost()
	assert.Contains(t, mockLogger.fatalFMsg, "FLYTE_API=\":/invalid url\" is not valid URL: ")
}

func TestHipchatAuthToken(t *testing.T) {

	os.Setenv("HIPCHAT_TOKENS", "abc")
	defer func() { os.Unsetenv("HIPCHAT_TOKENS") }()

	tokens := HipchatAuthTokens()
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "abc", tokens[0])
}

func TestHipchatAuthTokens(t *testing.T) {

	os.Setenv("HIPCHAT_TOKENS", "abc,  def , xyz,123,456")
	defer func() { os.Unsetenv("HIPCHAT_TOKENS") }()

	tokens := HipchatAuthTokens()
	assert.Equal(t, 5, len(tokens))
	assert.Equal(t, "abc", tokens[0])
	assert.Equal(t, "def", tokens[1])
	assert.Equal(t, "xyz", tokens[2])
	assert.Equal(t, "123", tokens[3])
	assert.Equal(t, "456", tokens[4])
}

func TestHipchatAuthTokensNotSet(t *testing.T) {

	mockLogger := NewMockLogger()
	defer func() { mockLogger.rollback() }()

	HipchatAuthTokens()
	assert.Equal(t, "env=HIPCHAT_TOKENS not set", mockLogger.fatalFMsg)
}

func TestDefaultRoom(t *testing.T) {
	assert.Equal(t, "", DefaultRoom())
}

func TestDefaultRoomSet(t *testing.T) {

	os.Setenv("DEFAULT_JOIN_ROOM", "abc")
	defer func() { os.Unsetenv("DEFAULT_JOIN_ROOM") }()

	assert.Equal(t, "abc", DefaultRoom())
}

func TestBkpDirDefault(t *testing.T) {
	assert.Equal(t, "", BkpDir())
}

func TestBkpDir(t *testing.T) {

	os.Setenv("BKP_DIR", "/tmp/hipchat-pack")
	defer func() { os.Unsetenv("BKP_DIR") }()

	assert.Equal(t, "/tmp/hipchat-pack", BkpDir())
}

type MockLogger struct {
	prevLogger func(string, ...interface{})
	fatalFMsg  string
}

func NewMockLogger() *MockLogger {

	l := &MockLogger{}
	l.prevLogger = logger.Fatalf
	logger.Fatalf = l.FatalF
	return l
}

func (l *MockLogger) FatalF(format string, a ...interface{}) {
	l.fatalFMsg = fmt.Sprintf(format, a...)
}

func (l *MockLogger) rollback() {
	logger.Fatalf = l.prevLogger
}
