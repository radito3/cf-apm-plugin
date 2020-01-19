package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type OperationMonitor struct {
	OperationId string
	Client      *HttpClient
}

func (monitor *OperationMonitor) monitorOperation() {
	for {
		statusCode, messages, err := monitor.pollForMessages()
		if err != nil {
			fmt.Println(err)
			break
		}

		switch statusCode {
		case 102:
			for _, msg := range messages {
				fmt.Println(msg)
			}
			time.Sleep(2 * time.Second)
		case 201:
			fmt.Println("Application updated")
			return
		case 300:
			fmt.Printf("Operation is in validation phase\n" +
				"Use \"cf bg-upload %s --continue\" to switch to new app\n", monitor.OperationId)
			return
		default:
			fmt.Println("Unknown error")
			return
		}
	}
}

func (monitor *OperationMonitor) pollForMessages() (sc int, msgs []string, err error) {
	request := HttpRequest{
		Method: "GET",
		Url:    monitor.Client.getBaseUrl() + "messages/" + monitor.OperationId,
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
