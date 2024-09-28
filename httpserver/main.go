package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Found / request\n")
	_, err := io.WriteString(w, "Root response from Webserver\n")
	if err != nil {
		return
	}
}

func getHelloMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got /hello message\n")
	_, err := io.WriteString(w, "Hello from HTTP server!\n")
	if err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHelloMessage)
	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
