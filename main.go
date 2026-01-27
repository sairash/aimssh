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
	flag.BoolVar(&runAsSSH, "ssh", false, "run as SSH server")
	flag.Parse()

	if runAsSSH {
		cfg := app.DefaultServerConfig()
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
