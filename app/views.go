package app

import (
	"fmt"
	"time"

	"aimssh/helper"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// View renders the current state of the application
func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var view string
	switch m.State {
	case LogoView:
		view = fmt.Sprintf("\n%s  \n\n\n            %s%s\n\n                  Loading...\n",
			TitleStyle.SetString(Logo).Render(),
			GreenColor.Bold(true).Render("ssh"),
			SelectedItemStyle.Underline(true).PaddingLeft(2).Render("aim.ftp.sh"))

	case InputView:
		view = fmt.Sprintf(
			"%s \n\n%s",
			TitleStyle.Render(),
			PaddingLeftStyle.Render(
				fmt.Sprintf("%s\n%s",
					HeightStyle.Render(
						fmt.Sprintf("%s \n\n%s\n\n%s \n\n%s",
							ListTitleStyle.Render("Time in minute: *"),
							m.Input.View(),
							BrownColor.Render("Session:"),
							m.WorkingOn.View(),
						),
					),
					m.helpView())))

	case ListView:
		view = fmt.Sprintf(
			"%s \n\n%s%s%s",
			TitleStyle.Render(),
			m.List.View(),
			DotStyle,
			m.List.Help.ShortHelpView([]key.Binding{m.Keymap.Back}),
		)

	case TimerView:
		text := "Session Ended, Press (r) or (n)"
		if !m.TimedOut {
			text = formatTime(int(m.Timer.Timeout.Minutes())) + " : " + formatTime(int(m.Timer.Timeout.Seconds())%60)
		}
		view = fmt.Sprintf(
			"%s \n\n %s",
			TitleStyle.Render(),
			PaddingLeftStyle.Render(fmt.Sprintf("%s\n%s%s\n%s",
				helper.Center(text, AppWidth-10),
				m.DrawASCII(m.Minute, m.Timer.Timeout),
				helper.Center("Session: "+m.Session, AppWidth-8),
				m.helpView())))
	}

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, AppStyle.Render(view))
}

// helpView returns the help text for the current state
func (m Model) helpView() string {
	switch m.State {
	case InputView:
		return "\n" + m.Help.ShortHelpView([]key.Binding{
			m.Keymap.Enter,
			m.Keymap.Back,
			m.Keymap.CtrlC,
		})
	}

	return "\n" + m.Help.ShortHelpView([]key.Binding{
		m.Keymap.Start,
		m.Keymap.Stop,
		m.Keymap.Reset,
		m.Keymap.Quit,
		m.Keymap.New,
	})
}

// DrawASCII renders the ASCII art animation based on timer progress
func (m Model) DrawASCII(total, remaining time.Duration) string {
	n := m.AsciiArt.Height()
	if n == 0 {
		return ""
	}
	if m.SelectedItem != "Flow" {
		return m.AsciiArt.NextAndString(int(percentageDifference(total, remaining)))
	}
	return "\n" + GreenColor.Render(m.AsciiArt.NextAndString(0)) + "\n"
}

// formatTime formats a number as a two-digit string with leading zero
func formatTime(n int) string {
	var b [2]byte
	if n < 10 {
		b[0], b[1] = '0', byte(n)+'0'
	} else {
		b[0], b[1] = byte(n/10)+'0', byte(n%10)+'0'
	}
	return string(b[:])
}

// percentageDifference calculates the percentage of time elapsed
func percentageDifference(total, remaining time.Duration) float64 {
	if total == 0 && remaining == 0 {
		return 0.0
	}
	return ((total.Seconds() - remaining.Seconds()) / total.Seconds()) * 100
}
