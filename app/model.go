package app

import (
	"fmt"
	"io"
	"strings"
	"time"

	"aimssh/ascii_generator"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

// SessionState represents the current view/state of the application
type SessionState int

const (
	InputView SessionState = iota
	LogoView
	ListView
	TimerView
)

// Item represents a selectable item in the list
type Item string

// FilterValue implements list.Item interface
func (i Item) FilterValue() string { return "" }

// ItemDelegate handles rendering of list items
type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// Model represents the application state
type Model struct {
	State        SessionState
	Input        textinput.Model
	WorkingOn    textinput.Model
	List         list.Model
	Minute       time.Duration
	SelectedItem string
	Timer        timer.Model
	LoadingTimer timer.Model
	Keymap       Keymap
	Help         help.Model
	Err          error
	Width        int
	Height       int
	Session      string
	AsciiArt     ascii_generator.AsciiArt
	Quitting     bool
	TimedOut     bool
	SessionSSH   ssh.Session
	RunAsSSH     bool
}

// NewModel creates and returns a new Model with default values
func NewModel(session interface{}, runAsSSH bool) Model {
	ti := textinput.New()
	ti.PlaceholderStyle = lipgloss.NewStyle().Faint(true)
	ti.Placeholder = "10"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30
	ti.Prompt = "- "

	woI := textinput.New()
	woI.PlaceholderStyle = lipgloss.NewStyle().Faint(true)
	woI.Placeholder = "Work"
	woI.CharLimit = 50
	woI.Width = 30
	woI.Prompt = "- "

	items := []list.Item{
		Item("Tree"),
		Item("Flow"),
		Item("Coffee"),
	}

	l := list.New(items, ItemDelegate{}, 30, 11)
	l.Title = "Select a visual option: "
	l.SetShowStatusBar(false)
	l.Styles.Title = ListTitleStyle
	l.SetHeight(23)

	m := Model{
		State:        LogoView,
		Input:        ti,
		LoadingTimer: timer.NewWithInterval(800*time.Millisecond, time.Millisecond),
		WorkingOn:    woI,
		List:         l,
		Err:          nil,
		Session:      "Work",
		TimedOut:     false,
		Keymap:       NewKeymap(),
		Help:         help.New(),
		RunAsSSH:     runAsSSH,
	}

	switch sess := session.(type) {
	case ssh.Session:
		m.SessionSSH = sess
	}

	return m
}

// Init implements tea.Model interface
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.LoadingTimer.Init(), m.LoadingTimer.Start())
}

// GenerateASCII creates the appropriate ASCII art based on selected item
func (m Model) GenerateASCII() ascii_generator.AsciiArt {
	switch m.SelectedItem {
	case "Coffee":
		return ascii_generator.GenerateCoffee()
	case "Tree":
		return ascii_generator.GenerateTree(40, 18)
	default:
		return ascii_generator.GenerateRow(40, 17)
	}
}
