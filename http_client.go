package main

import (
	"bytes"
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"io"
	"mime/multipart"
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

func createRequestBodyWithFile(fileName string) (io.Reader, map[string]string) {
	file, err := os.Open(fileName)
	if err != nil{
		fmt.Println(err)
		return nil, nil
	}

	defer func() {
		err1 := file.Close()
		if err1 != nil {
			fmt.Println(err1)
		}
	}()

	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	fileWriter, err := multiPartWriter.CreateFormFile("file_field", fileName)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	err2 := multiPartWriter.Close()
	if err2 != nil {
		fmt.Println(err2)
		return nil, nil
	}

	params := make(map[string]string, 2)
	params["Content-Length"] = string(info.Size())
	params["Content-Type"] = multiPartWriter.FormDataContentType()

	return &requestBody, params
}

func (c *HttpClient) getBaseUrl() string {
	domain := strings.Join(strings.Split(c.Api, ".")[2:], ".")

	return fmt.Sprintf("https://%s.cfapps.%s/api/v1/", AppUrl, domain)
}

type HttpRequest struct {
	Method  string
	Url     string
	Token   string
	Body    io.Reader
	Params  map[string]string
}

func httpCall(request HttpRequest) (*http.Response, error) {
	req, err := http.NewRequest(request.Method, request.Url, request.Body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", request.Token)
	if request.Params != nil {
		for key, val := range request.Params {
			req.Header.Set(key, val)
		}
	}

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
