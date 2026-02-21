// AimSSH is a terminal-based Pomodoro timer application that can run
// locally as a TUI or as an SSH server for remote access.
package main

import (
	"flag"
	"fmt"
	"os"

	"aimssh/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var runAsSSH bool
	var host string
	var port string

	flag.BoolVar(&runAsSSH, "ssh", false, "run as SSH server")
	flag.StringVar(&host, "host", "0.0.0.0", "host to listen on")
	flag.StringVar(&port, "port", "13234", "port to listen on")

	flag.Parse()

	if runAsSSH {
		cfg := app.DefaultServerConfig()
		cfg.Host = host
		cfg.Port = port

		app.RunServer(cfg)
		return
	}

	p := tea.NewProgram(app.NewModel(nil, false), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if model, ok := m.(app.Model); ok && model.Quitting {
		if model.Err != nil {
			fmt.Printf("Error Occurred: %s \n", model.Err.Error())
			return
		}
		fmt.Println(app.EndInfo)
	}
}
