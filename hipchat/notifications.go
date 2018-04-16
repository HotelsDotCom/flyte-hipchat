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

var startUpNotification = Notification{
	Message:       "HipChat start up...",
	MessageFormat: "text",
	Color:         "green",
	From:          "flyte-hipchat",
}

var shutDownNotification = Notification{
	Message:       "HipChat shutting down...",
	MessageFormat: "text",
	Color:         "red",
	From:          "flyte-hipchat",
}

var joinNotification = Notification{
	Message:       "Hello! I've joined this room...",
	MessageFormat: "text",
	Color:         "green",
	From:          "flyte-hipchat",
}

var leaveNotification = Notification{
	Message:       "I'm leaving now, bye!",
	MessageFormat: "text",
	Color:         "red",
	From:          "flyte-hipchat",
}
