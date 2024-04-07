package tcp

import (
	"context"
	"fmt"
	"go_mini_redis/interface/tcp"
	"go_mini_redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Config stores tcp server properties
type Config struct {
	Address string `yaml:"address"`
}

// ListenAndServeWithSignal binds port and handle requests, blocking until receive stop signal
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	listerner, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	// close channel
	closeChan := make(chan struct{})

	// signal channel, recieve from OS signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	// If recieve OS stop signals, sent it to closeChan
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listerner, handler, closeChan)

	return nil
}

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {

	// Close listener and handler with recieving closeChan
	go func() {
		<-closeChan
		_ = listener.Close()
		_ = handler.Close()
	}()

	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()

	var waitDone sync.WaitGroup
	for true {
		conn, err := listener.Accept()

		logger.Info("link accepted")
		if err != nil {
			break
		}
		waitDone.Add(1)

		go func() {
			defer func() {
				waitDone.Done()
			}()
			// Forward request to handler
			handler.Handle(ctx, conn)
		}()
	}

	waitDone.Wait()
}
