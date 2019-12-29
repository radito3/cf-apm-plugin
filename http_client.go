package main

import (
	"bufio"
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type HttpClient struct {
	Api   string
	Token string
}

const AppUrl string = "blue-green-uploader"

func createRequestBodyWithFile(fileName string) io.Reader {
	f, err := os.Open(fileName)
	if err != nil{
		fmt.Println(err)
		return nil
	}

	return bufio.NewReader(f)
}

func (c *HttpClient) getBaseUrl() string {
	domain := strings.Join(strings.Split(c.Api, ".")[2:], ".")

	return fmt.Sprintf("https://%s.cfapps.%s/Api/v1/", AppUrl, domain)
}

type HttpRequest struct {
	Method string
	Url    string
	Token  string
	Body   io.Reader
}

func httpCall(request HttpRequest) (*http.Response, error) {
	req, err := http.NewRequest(request.Method, request.Url, request.Body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", request.Token)

	client := http.Client{Timeout: 5 * time.Minute}

	return client.Do(req)
}

func createHttpClient(con plugin.CliConnection) (*HttpClient, error) {
	api, err := con.ApiEndpoint()
	if err != nil {
		return nil, err
	}

	token, err := con.AccessToken()
	if err != nil {
		return nil, err
	}

	return &HttpClient{api, token}, nil
}
