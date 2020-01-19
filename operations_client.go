package main

import "fmt"

type OperationsClient struct {
	HttpClient  HttpClient
	OperationId string
}

func (o *OperationsClient) uploadApp(appName, fileName string) bool {
	body, params := createRequestBodyWithFile(fileName)
	request := HttpRequest{
		Method: "POST",
		Url:    o.HttpClient.getBaseUrl() + "upload/" + appName,
		Token:  o.HttpClient.Token,
		Body:   body,
		Params: params,
	}

	resp, err := httpCall(request)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if resp.StatusCode != 200 {
		fmt.Println("Server error")
		return false
	}

	operationId := resp.Header["Location"][0]
	o.OperationId = operationId
	return true
}

func (o *OperationsClient) continueAppUpload(operationId string) bool {
	request := HttpRequest{
		Method: "PUT",
		Url:    o.HttpClient.getBaseUrl() + "resume/" + operationId,
		Token:  o.HttpClient.Token,
	}

	resp, err := httpCall(request)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if resp.StatusCode != 200 {
		fmt.Println("Server error")
		return false
	}
	return true
}
