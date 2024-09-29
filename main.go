package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"http-in-go/httpserver"
	"io"
	"net/http"
)

func main() {
	httpserver.CreateHTTPServer("8080")

	requestUrl := fmt.Sprintf("http://localhost:%s", "8080")
	res, err := http.Get(requestUrl)
	if err != nil {
		logrus.Errorf("Error making http request: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	logrus.Info("Client: got response!\n")
	logrus.Infof("Client: response status code: %v\n", res.StatusCode)

}
