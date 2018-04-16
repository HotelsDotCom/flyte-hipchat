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

package client

import (
	"fmt"
	"testing"
)

func Test_SomeFailureSucceedsInTheEnd(t *testing.T) {
	var attempt = 0
	err := do(func() (err error) {
		attempt++
		if attempt > 2 {
			return nil
		} else {
			return fmt.Errorf("agh - something went wrong")
		}
	})
	if err != nil {
		t.Errorf("Got an error still %s", err)
	}
}

func Test_ExceedMaxRetriesError(t *testing.T) {
	var actualRetries = 0
	err := do(func() (err error) {
		actualRetries++
		return fmt.Errorf("agh - something wrong")
	})
	if err == nil {
		t.Errorf("No error returned %s", err)
	}
	if actualRetries != 10 {
		t.Errorf("Expected 10 retries, got %d", actualRetries)
	}
}

func Test_NoRetriesNeeded(t *testing.T) {
	err := do(func() (err error) {
		return nil
	})
	if err != nil {
		t.Errorf("Got an error still %s", err)
	}
}
