package main

import "github.com/ch55secake/http-in-go/httpserver"

func main() {
	err := httpserver.CreateHTTPServer("8080")
	if err != nil {
		return
	}
}
