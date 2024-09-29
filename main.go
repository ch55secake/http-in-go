package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"http-in-go/httpserver"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	httpserver.CreateHTTPServer("8080", wg)

	select {
	case <-time.After(time.Second * 10):
	}

	ping()
}

func ping() {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	resp, err := c.Get("http://localhost:8080")
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	_, err = ioutil.ReadAll(resp.Body)

	logrus.Info("Client: got response!\n")
	logrus.Infof("Client: response status code: %v\n", resp.StatusCode)
}
