package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// logDefaultMessage /**
// Create little ASCII art when server begins running
func logDefaultMessage() {
	startupMessage := ` 
 	 _________                                
 	/   _____/ ______________  __ ___________ 
 	\_____  \_/ __ \_  __ \  \/ // __ \_  __ \
	 /        \  ___/|  | \/\   /\  ___/|  | \/
	/_______  /\___  >__|    \_/  \___  >__|   
      	    \/     \/                 \/       `
	fmt.Println(startupMessage)
}

// CreateHTTPServer /**
// Creates a default httpServer with its own http.ServeMux, will only server one / endpoint as part of this
// Port can be provided as wanted by the user of this method
func CreateHTTPServer(p string, wg *sync.WaitGroup) {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: "15:04:05", FullTimestamp: true})
	logrus.SetOutput(colorable.NewColorableStdout())
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Server handling %v request at /", r.Method)
	})
	s := http.Server{
		Addr:    fmt.Sprintf(":%v", p),
		Handler: mux,
	}

	go func() {
		defer wg.Done()

		logDefaultMessage()
		logrus.Infof("Server has started running at http://localhost%v", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logrus.Errorf("Error running http server: %s\n", err)
			}
		}
	}()
	gracefulShutdown(&s)
}

// CreateHTTPServerWithMux /**
// Creates HTTPServer, but allows user to pass their own http.ServeMux which could serve any number of endpoints
func CreateHTTPServerWithMux(p string, mux *http.ServeMux) error {
	logDefaultMessage()
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: "15:04:05", FullTimestamp: true})
	logrus.SetOutput(colorable.NewColorableStdout())
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", p),
		Handler: mux,
	}

	logDefaultMessage()
	logrus.Infof("Server has started running at http://localhost%v", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("Error running http server: %s\n", err)
		}
	}
	return nil
}

// CreateHttpServerWithMuxAndContext /**
// Creates HttpServer, but allows user to pass their own http.ServeMux and context.Context, allow user to provide context
// and allowing user to create their own endpoints
func CreateHttpServerWithMuxAndContext(p string, mux *http.ServeMux, ctx context.Context, addr string) error {
	logDefaultMessage()
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: "15:04:05", FullTimestamp: true})
	logrus.SetOutput(colorable.NewColorableStdout())

	ctx, cancelCtx := context.WithCancel(context.Background())
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", p),
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			ctx = context.WithValue(ctx, addr, listener.Addr().String())
			return ctx
		},
	}

	go func() {
		logrus.Infof("Server has started running at http://localhost%v", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("Error running http server: %s\n", err)
		}
		logrus.Info("Stopped serving new connections.")
	}()

	cancelCtx()
	return nil
}

func gracefulShutdown(s *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	shutDownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := s.Shutdown(shutDownCtx); err != nil {
		logrus.Errorf("Error shutting down server: %v", err)
	}
	logrus.Infof("Server has stopped serving new connections.")
}
