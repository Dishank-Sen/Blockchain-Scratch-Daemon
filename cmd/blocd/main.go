package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/internal/daemon"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
)

func init() {
	logPath := defaultLogPath()

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		logger.Error(fmt.Sprintf("failed to create log dir: %v", err))
		return
	}

	logFile, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to open log file: %v", err))
		return
	}

	log.SetOutput(logFile)
}

func defaultLogPath() string {
	if runtime.GOOS == "windows" {
		base := os.Getenv("LOCALAPPDATA")
		if base == "" {
			base = os.TempDir()
		}
		return filepath.Join(base, "blocd", "blocd.log")
	}

	// Linux / macOS
	return "/tmp/blocd.log"
}


func main(){
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	addr := ":4001"
	daemon, err := daemon.NewDaemon(ctx, addr)
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