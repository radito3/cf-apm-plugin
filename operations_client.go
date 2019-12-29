package bgUploaderPlugin

import "fmt"

type OperationsClient struct {
	httpClient HttpClient
}

func (o *OperationsClient) uploadApp(appName, fileName string) {
	request := HttpRequest{
		Method: "POST",
		Url:    o.httpClient.getBaseUrl() + "/upload/" + appName,
		Token:  o.httpClient.token,
		Body:   createRequestBodyWithFile(fileName),
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

func (o *OperationsClient) continueAppUpload(appName string) {
	request := HttpRequest{
		Method: "PUT",
		Url:    o.httpClient.getBaseUrl() + "/upload/" + appName,
		Token:  o.httpClient.token,
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
