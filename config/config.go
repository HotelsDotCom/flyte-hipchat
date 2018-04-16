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
	"net/url"
	"os"
	"github.com/HotelsDotCom/go-logger"
	"strings"
)

func ApiHost() *url.URL {

	hostEnv := getEnv("FLYTE_API", true)
	host, err := url.Parse(hostEnv)
	if err != nil {
		logger.Fatalf("FLYTE_API=%q is not valid URL: %v", hostEnv, err)
	}
	return host
}

func HipchatAuthTokens() []string {

	tokensEnv := getEnv("HIPCHAT_TOKENS", true)
	tokens := []string{}
	for _, t := range strings.Split(tokensEnv, ",") {
		tokens = append(tokens, strings.TrimSpace(t))
	}
	return tokens
}

func DefaultRoom() string {
	return getEnv("DEFAULT_JOIN_ROOM", false)
}

func BkpDir() string {
	return getEnv("BKP_DIR", false)
}

func getEnv(key string, required bool) string {

	v := os.Getenv(key)
	if required && v == "" {
		logger.Fatalf("env=%s not set", key)
	}
	return v
}
