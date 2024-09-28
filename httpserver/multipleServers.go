package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func getBase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")

	fmt.Printf("%s: got / request! first(%t)=%s\n", ctx.Value(keyServerAddr), isFirst, first)
	_, err := io.WriteString(w, "Base banana!\n")
	if err != nil {
		return
	}
}

func getWelcomeMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	_, err := io.WriteString(w, "Hello from multiple servers!\n")
	if err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getBase)
	mux.HandleFunc("/hello", getWelcomeMessage)

	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	serverTwo := &http.Server{
		Addr:    ":4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := serverOne.ListenAndServe()
		if err != http.ErrServerClosed {
			fmt.Printf("serverOne closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server one: %s\n ", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()

	go func() {
		err := serverTwo.ListenAndServe()
		if err != http.ErrServerClosed {
			fmt.Printf("serverTwo closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server two: %s\n ", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()

}
