package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mqdvi-dp/go-common/factory"
	"github.com/mqdvi-dp/go-common/logger"
)

type Server struct {
	service factory.ServiceFactory
}

// New service server
func New(service factory.ServiceFactory) *Server {
	return &Server{service: service}
}

func (s *Server) Run() {
	if len(s.service.GetApplications()) < 1 {
		logger.Log.Fatal("No server/worker running")
	}

	errServe := make(chan error, len(s.service.GetApplications()))
	for _, app := range s.service.GetApplications() {
		go func(srv factory.AppServerFactory) {
			defer func() {
				if r := recover(); r != nil {
					errServe <- fmt.Errorf("%s", r)
				}
			}()

			// run server
			srv.Serve()
		}(app)
	}

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt)
	signal.Notify(quitSignal, syscall.SIGTERM)

	fmt.Println()
	fmt.Printf("Application \x1b[32;1m%s\x1b[0m ready to run\n", s.service.Name())

	select {
	case e := <-errServe:
		panic(e)
	case <-quitSignal:
		s.shutdown(quitSignal)
	}
}

// graceful shutdown all server, panic if there is still a process running when the requests exceed given timeout in context
func (s *Server) shutdown(forceShutdown chan os.Signal) {
	fmt.Println("\x1b[34;1mGracefully shutdown... (press Ctrl+C again to force)\x1b[0m")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for _, server := range s.service.GetApplications() {
			server.Shutdown(ctx)
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
		fmt.Printf("\x1b[32;1mSuccess shutdown all server & worker at %s\x1b[0m\n", s.service.Name())
	case <-forceShutdown:
		fmt.Println("\x1b[31;1mForce shutdown server & worker\x1b[0m")
		cancel()
	case <-ctx.Done():
		fmt.Println("\x1b[31;1mContext timeout\x1b[0m")
	}
}
