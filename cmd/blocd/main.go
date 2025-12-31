package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/daemon"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

func init() {
	logFile, err := os.OpenFile(
		"/tmp/blocd.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		// If logging fails, fall back to stderr
		logger.Error(fmt.Sprintf("failed to open log file: %v", err))
		return
	}

	log.SetOutput(logFile)
}

func main(){
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	
	daemon, err := daemon.NewDaemon(ctx)
	if err != nil{
		logger.Error(fmt.Sprintf("new daemon error: %v", err))
		stop()
	}
	logger.Info("Daemon started..")

	// blocking code
	if err := daemon.Run(); err != nil{
		logger.Error(fmt.Sprintf("error while running daemon: %v", err))
		stop()
	}
}