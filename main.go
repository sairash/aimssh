// AimSSH is a terminal-based Pomodoro timer application that can run
// locally as a TUI or as an SSH server for remote access.
package main

import (
	"fmt"
	"os"

	"aimssh/app"

	tea "github.com/charmbracelet/bubbletea"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.1"

	sshHost string
	sshPort string
)

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "aimssh",
	Short: "AimSSH is a terminal-based Pomodoro timer",
	Run: func(cmd *cobra.Command, args []string) {
		runTUI()
	},
}

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Run AimSSH as an SSH server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := app.DefaultServerConfig()
		cfg.Host = sshHost
		cfg.Port = sshPort

		app.RunServer(cfg)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of AimSSH",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var genManCmd = &cobra.Command{
	Use:   "genman",
	Short: "Use this command to generate the man page",
	Run: func(cmd *cobra.Command, args []string) {
		manPage, err := mcobra.NewManPage(1, rootCmd)
		if err != nil {
			panic(err)
		}

		manPage = manPage.WithSection("Copyright", "(C) 2026 Sairash Sharma Gautam. \n"+"Released under MIT license.")
		fmt.Println(manPage.Build(roff.NewDocument()))
	},
}

func init() {
	sshCmd.Flags().StringVarP(
		&sshHost,
		"host",
		"u",
		"0.0.0.0",
		"Host address to run the SSH server on",
	)

	sshCmd.Flags().StringVarP(
		&sshPort,
		"port",
		"p",
		"13234",
		"Port to run the SSH server on",
	)

	rootCmd.AddCommand(sshCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(genManCmd)
}

func runTUI() {
	p := tea.NewProgram(
		app.NewModel(nil, false),
		tea.WithAltScreen(),
	)

	m, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if model, ok := m.(app.Model); ok && model.Quitting {
		if model.Err != nil {
			fmt.Printf("Error Occurred: %s\n", model.Err.Error())
			return
		}
		fmt.Println(app.EndInfo)
	}
}
