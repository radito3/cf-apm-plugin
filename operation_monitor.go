package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type OperationMonitor struct {
	Client      *HttpClient
}

func (monitor *OperationMonitor) monitorOperation(operationId string) {
	for {
		statusCode, messages, err := monitor.pollForMessages(operationId)
		if err != nil {
			fmt.Println(err)
			break
		}

		switch statusCode {
		case 200:
			for _, msg := range messages {
				fmt.Println(msg)
			}
			time.Sleep(3 * time.Second)
		case 201:
			fmt.Println("Application updated")
			return
		default:
			fmt.Println("Unknown error")
			return
		}
	}
}

func (monitor *OperationMonitor) pollForMessages(operationId string) (sc int, msgs []string, err error) {
	request := HttpRequest{
		Method: "GET",
		Url:    monitor.Client.getBaseUrl() + "messages/" + operationId,
		Token:  monitor.Client.Token,
	}
	sc = -1

	response, err1 := httpCall(request)
	if err1 != nil {
		err = err1
		return
	}

	defer func() {
		e := response.Body.Close()
		if e != nil {
			err = e
		}
	}()

	bytes, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		err = err1
		return
	}

	var messages []string

	err2 := json.Unmarshal(bytes, &messages)
	if err2 != nil {
		err = err2
		return
	}

	return response.StatusCode, messages, nil
}
