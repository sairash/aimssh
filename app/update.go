package app

import (
	"fmt"
	"strconv"
	"time"

	"aimssh/helper"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all state transitions based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.Quitting = true
			return m, tea.Batch(tea.Quit, helper.BeepCmd())

		case tea.KeyCtrlB:
			switch m.State {
			case InputView:
				if m.WorkingOn.Focused() {
					m.WorkingOn.Blur()
					m.Input.Focus()
				}
			case ListView:
				m.State = InputView
				m.Input.Blur()
				m.WorkingOn.Focus()
			case TimerView:
				m.State = ListView
			}
		}
	}

	switch m.State {
	case LogoView:
		return updateLogo(msg, m)
	case InputView:
		return updateInput(msg, m)
	case ListView:
		return updateList(msg, m)
	case TimerView:
		return updateTimer(msg, m)
	default:
		return m, nil
	}
}

// updateLogo handles updates for the logo/splash screen view
func updateLogo(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.LoadingTimer, cmd = m.LoadingTimer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.LoadingTimer.Stop()
		m.State = InputView
		return m, textinput.Blink
	}
	return m, nil
}

// updateInput handles updates for the input view (time and session name)
func updateInput(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Input.Focused() {
				if m.Input.Value() != "" {
					min, err := strconv.Atoi(m.Input.Value())
					if err != nil {
						m.Err = err
						m.Quitting = true
						return m, tea.Quit
					}
					m.Minute = time.Duration(min) * time.Minute
					m.Input.Blur()
					return m, m.WorkingOn.Focus()
				}
			} else {
				if m.WorkingOn.Value() != "" {
					m.Session = m.WorkingOn.Value()
				}
				m.State = ListView
			}
		}
	}

	if m.Input.Focused() {
		m.Input, cmd = m.Input.Update(msg)
	} else {
		m.WorkingOn, cmd = m.WorkingOn.Update(msg)
	}

	return m, cmd
}

// updateList handles updates for the visual option selection view
func updateList(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if selected, ok := m.List.SelectedItem().(Item); ok {
				m.SelectedItem = string(selected)
				m.Timer = timer.NewWithInterval(m.Minute, time.Millisecond)
				m.State = TimerView
				m.AsciiArt = m.GenerateASCII()
				m.Keymap.Start.SetEnabled(false)

				body := fmt.Sprintf("Timer set for %d minutes.", int(m.Timer.Timeout.Minutes()))
				SendNotification(m.SessionSSH, helper.TimerStartedTitle, body, m.RunAsSSH)

				return m, m.Timer.Init()
			}
		}
	}

	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

// updateTimer handles updates for the timer view
func updateTimer(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.TimedOut = true
		body := "The timer has ended."
		SendNotification(m.SessionSSH, helper.TimerEndedTitle, body, m.RunAsSSH)
		return m, helper.BeepCmd()

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		m.Keymap.Stop.SetEnabled(m.Timer.Running())
		m.Keymap.Start.SetEnabled(!m.Timer.Running())
		return m, tea.Batch(cmd, helper.BeepCmd())

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.Quit):
			m.TimedOut = false
			m.Quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.Keymap.Reset):
			m.TimedOut = false
			m.AsciiArt = m.GenerateASCII()
			m.Timer.Timeout = m.Minute

			body := fmt.Sprintf("Timer set for %d minutes.", int(m.Timer.Timeout.Minutes()))
			SendNotification(m.SessionSSH, helper.TimerRestartedTitle, body, m.RunAsSSH)

			return m, m.Timer.Start()

		case key.Matches(msg, m.Keymap.New):
			m.TimedOut = false
			m.State = InputView
			m.WorkingOn.Blur()
			m.Input.Focus()
			m.Timer.Stop()
			return m, nil

		case key.Matches(msg, m.Keymap.Start, m.Keymap.Stop):
			return m, m.Timer.Toggle()
		}
	}

	return m, nil
}
