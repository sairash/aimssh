package app

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

// ServerConfig holds configuration for the SSH server
type ServerConfig struct {
	Host       string
	Port       string
	HostKeyPath string
}

// DefaultServerConfig returns the default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:       "localhost",
		Port:       "13234",
		HostKeyPath: ".ssh/id_ed25519",
	}
}

// RunServer starts the SSH server with the given configuration
func RunServer(cfg ServerConfig) {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(cfg.Host, cfg.Port)),
		wish.WithHostKeyPath(cfg.HostKeyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
			// Middleware to print end info after Bubbletea quits
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					wish.Println(sess, EndInfo)
					next(sess)
				}
			},
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", cfg.Host, "port", cfg.Port)

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

// teaHandler creates a new bubbletea program for each SSH session
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	s.Context().SetValue("operating_system", detectOS(s))
	return NewModel(s, true), []tea.ProgramOption{tea.WithAltScreen()}
}

// detectOS attempts to detect the client's operating system
func detectOS(s ssh.Session) string {
	if wish.Command(s, "paplay", "--help").Run() == nil {
		return "linux"
	}
	return "windows"
}
