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

package main

import (
	"net/url"
	"os"
	"os/signal"
	api "github.com/HotelsDotCom/flyte-client/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/flyte-hipchat/bkp"
	"github.com/HotelsDotCom/flyte-hipchat/client"
	"github.com/HotelsDotCom/flyte-hipchat/command"
	"github.com/HotelsDotCom/flyte-hipchat/config"
	"github.com/HotelsDotCom/flyte-hipchat/event"
	"github.com/HotelsDotCom/flyte-hipchat/hipchat"
	"github.com/HotelsDotCom/go-logger"
	"syscall"
	"time"
)

func main() {

	messages := make(chan hipchat.Message)
	hc := initHipchat(messages)

	p := flyte.NewPack(getPackDef(hc), api.NewClient(config.ApiHost(), 10*time.Second))
	p.Start()

	event.HandleReceivedMessages(p, messages)
	logger.Infof("joined rooms=%v", hc.JoinedRoomIds())

	// block until we get an exit causing signal
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		logger.Info("received interrupt, shutting down...")
		hc.Shutdown()
		logger.Info("shut down")
	}
}

func initHipchat(messages chan hipchat.Message) hipchat.Hipchat {

	hcClient := client.NewHipChatClient(config.HipchatAuthTokens())

	bkpDir := config.BkpDir()
	if bkpDir == "" {
		bkpDir = bkp.CreateDefaultBkpDir()
	}

	bkpFile := bkp.CreateBkpFile(bkpDir, "rooms.json")
	hc, err := hipchat.NewHipchat(bkpFile, hcClient, messages)
	if err != nil {
		logger.Fatalf("cannot initialize pack: %v", err)
	}

	if defaultRoom := config.DefaultRoom(); defaultRoom != "" {
		hc.JoinRoom(defaultRoom)
	}

	if len(hc.JoinedRoomIds()) == 0 {
		logger.Fatal("pack did NOT join any rooms, provide DEFAULT_JOIN_ROOM env. var.")
	}
	return hc
}

func getPackDef(hc hipchat.Hipchat) flyte.PackDef {

	helpUrl, err := url.Parse("http://github.com/HotelsDotCom/flyte-hipchat/browse/README.md")
	if err != nil {
		logger.Fatal("invalid pack help url")
	}

	return flyte.PackDef{
		Name:    "HipChat",
		HelpURL: helpUrl,
		Commands: []flyte.Command{
			command.SendMessageCommand(hc),
			command.SendNotificationCommand(hc),
			command.BroadcastCommand(hc),
			command.JoinCommand(hc),
			command.LeaveCommand(hc),
		},
		EventDefs: []flyte.EventDef{
			{Name: "ReceivedMessage"},
		},
	}
}
