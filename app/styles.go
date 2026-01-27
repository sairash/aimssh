// Package app provides the core application logic for AimSSH,
// a terminal-based Pomodoro timer that can run locally or as an SSH server.
package app

import (
	"fmt"

	"aimssh/helper"

	"github.com/charmbracelet/lipgloss"
)

// Application constants
const (
	DotChar  = " • "
	AppWidth = 50
	Logo     = `    _    ___ __  __   ____ ____  _   _ 
   / \  |_ _|  \/  | / ___/ ___|| | | |
  / _ \  | || |\/| | \___ \___ \| |_| |
 / ___ \ | || |  | |  ___) |__) |  _  |
/_/   \_\___|_|  |_| |____/____/|_| |_|
                                              `
)

// Style definitions for the application UI
var (
	AppStyle          = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder(), true, true, true, true).Width(AppWidth)
	HeightStyle       = lipgloss.NewStyle().Height(21)
	PaddingLeftStyle  = lipgloss.NewStyle().PaddingLeft(2)
	TitleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#49beaa")).Bold(true).SetString(helper.Center(`<AIM SSH>`, AppWidth-4)).AlignHorizontal(lipgloss.Center)
	ListTitleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#bfedc1")).PaddingLeft(-10)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#CFF27E"))
	GreenColor        = lipgloss.NewStyle().Foreground(lipgloss.Color("#bfedc1")).PaddingLeft(2).Faint(true)
	DotStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(DotChar)
	BrownColor        = lipgloss.NewStyle().Foreground(lipgloss.Color("#967969"))
)

// GitLink is the styled GitHub repository link
var GitLink = GreenColor.Render("https://github.com/sairash/aimssh")

// EndInfo is the message displayed when the application exits
var EndInfo = fmt.Sprintf("\n Thanks for using %s! \n Give a star %s \n Made By     %s\n",
	lipgloss.NewStyle().Foreground(lipgloss.Color("#49beaa")).Bold(true).Render("<AIM SSH>"),
	GitLink,
	GreenColor.Render("https://sairashgautam.com.np/"))
