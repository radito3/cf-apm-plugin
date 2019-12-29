package main

import "fmt"

type OperationsClient struct {
	HttpClient  HttpClient
	OperationId string
}

func (o *OperationsClient) uploadApp(appName, fileName string) {
	request := HttpRequest{
		Method: "POST",
		Url:    o.HttpClient.getBaseUrl() + "upload/" + appName,
		Token:  o.HttpClient.Token,
		Body:   createRequestBodyWithFile(fileName),
	}

	resp, err := httpCall(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	operationId := resp.Header["Location"][0]
	o.OperationId = operationId

	if resp.StatusCode != 200 {
		fmt.Println("Server error")
	}
}

func (o *OperationsClient) continueAppUpload(operationId string) {
	request := HttpRequest{
		Method: "PUT",
		Url:    o.HttpClient.getBaseUrl() + "resume/" + operationId,
		Token:  o.HttpClient.Token,
	}

	resp, err := httpCall(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println("Server error")
	}
}
