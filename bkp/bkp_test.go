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

package bkp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"github.com/HotelsDotCom/go-logger"
	"testing"
)

func TestCreateDefaultBkpDir(t *testing.T) {

	// do NOT 'clean up' default bkp, so we don't delete existing room(s)
	assert.Contains(t, CreateDefaultBkpDir(), "/flyte-hipchat")
}

func TestCreateDefaultBkpFile(t *testing.T) {

	path := CreateBkpFile(CreateDefaultBkpDir(), "test-rooms.json")
	defer func() {
		if err := os.Remove(path); err != nil {
			logger.Errorf("cannot remove path %s: %v", path, err)
		}
	}()

	assert.Contains(t, path, "/flyte-hipchat/test-rooms.json")
}

func TestCreateBkpFile(t *testing.T) {

	bkpDir := filepath.Join(os.TempDir(), "flyte-test-hipchat")
	os.MkdirAll(bkpDir, 0755)

	defer func() {
		if err := os.RemoveAll(bkpDir); err != nil {
			logger.Errorf("cannot remove dir %s: %v", bkpDir, err)
		}
	}()

	path := CreateBkpFile(bkpDir, "rooms.json")
	assert.Contains(t, path, "/flyte-test-hipchat/rooms.json")
}

func TestCreateBkpFileNonExistingDir(t *testing.T) {

	mockLogger := NewMockLogger()
	defer func() { mockLogger.rollback() }()

	CreateBkpFile("/non-existing-dir/", "rooms.json")
	assert.Contains(t, mockLogger.fatalFMsg, "cannot write initial bkp file /non-existing-dir/rooms.json: ")
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
