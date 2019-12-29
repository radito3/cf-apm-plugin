package bgUploaderPlugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type OperationMonitor struct {
	AppName string
	Client  *HttpClient
}

func (monitor *OperationMonitor) monitorOperation() {
	var hasError = false
	for !hasError {
		request := HttpRequest{
			Method: "GET",
			Url:    monitor.Client.getBaseUrl() + "/messages/" + monitor.AppName,
			Token:  monitor.Client.token,
		}

		response, err := httpCall(request)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer func() {
			e := response.Body.Close()
			if e != nil {
				fmt.Println(e)
				hasError = true
			}
		}()

		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		var messages []string

		err2 := json.Unmarshal(bytes, &messages)
		if err2 != nil {
			fmt.Println(err2)
			return
		}

		switch response.StatusCode {
		case 102:
			if len(messages) > 0 {
				for msg := range messages {
					fmt.Println(msg)
				}
			}
			time.Sleep(2 * time.Second)
		case 201:
			fmt.Println("Application updated")
			return
		case 300:
			fmt.Printf("Operation is in validation phase\n" +
				"Use \"cf bg-upload %s --continue\" to switch to new app\n", monitor.AppName)
			return
		case 500:
			fmt.Println("Server error: " + messages[0])
			return
		default:
			fmt.Println("Unknown error")
			return
		}
	}
}
