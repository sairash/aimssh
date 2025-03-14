package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"zencli/ascii_generator"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	dotChar                = " • "
	app_width              = 50
	inputView sessionState = iota
	listView
	timerView
)

var (
	appStyle          = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder(), true, true, true, true).Width(app_width)
	heightThing       = lipgloss.NewStyle().Height(22)
	paddingleft       = lipgloss.NewStyle().PaddingLeft(2)
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#49beaa")).Bold(true).SetString(center("Zen Cli", app_width-4)).AlignHorizontal(lipgloss.Center)
	listTitleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#bfedc1")).PaddingLeft(-10)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#CFF27E"))

	skyColor    = lipgloss.NewStyle().Background(lipgloss.Color("#6495ED"))
	dotStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	state        sessionState
	input        textinput.Model
	list         list.Model
	minute       time.Duration
	selectedItem string
	timer        timer.Model
	keymap       keymap
	help         help.Model
	err          error
	width        int
	height       int
	asciiArt     ascii_generator.AsciiArt
	quitting     bool
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "10"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30
	ti.Prompt = "- "

	items := []list.Item{
		item("None"),
		item("Tree"),
		item("Tomato"),
		item("Coffee"),
		item("Carrot"),
	}

	l := list.New(items, itemDelegate{}, 30, 11)
	l.Title = "Select a visual option: "
	l.SetShowStatusBar(false)
	l.Styles.Title = listTitleStyle
	// l.SetShowHelp(false)
	l.SetHeight(24)
	// l.SetShowTitle(false)

	return model{
		state: inputView,
		input: ti,
		list:  l,
		err:   nil,
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "restart"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.quitting = true
			return m, tea.Quit
		}
	}

	switch m.state {
	case inputView:
		return updateInput(msg, m)
	case listView:
		return updateList(msg, m)
	case timerView:
		return updateTimer(msg, m)
	default:
		return m, nil
	}
}

func updateInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.input.Value() != "" {

				min, err := strconv.Atoi(m.input.Value())

				if err != nil {
					m.err = err
					m.quitting = true
					return m, tea.Quit
				}

				m.minute = time.Duration(min) * time.Minute
				m.state = listView
				return m, nil
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func updateList(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if selected, ok := m.list.SelectedItem().(item); ok {
				m.selectedItem = string(selected)
				m.timer = timer.NewWithInterval(m.minute, time.Millisecond)
				m.state = timerView
				m.asciiArt = ascii_generator.GenerateTree(40, 20)
				// +brownColor.Render(strings.Repeat("░", app_width-8))
				m.keymap.start.SetEnabled(false)
				return m, m.timer.Init()
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateTimer(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	// case timer.TimeoutMsg:
	// 	m.quitting = true
	// 	return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.asciiArt = ascii_generator.GenerateTree(40, 20)
			m.timer.Timeout = m.minute
			return m, m.timer.Start()
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m model) DrawAscii(a, b time.Duration) string {
	n := m.asciiArt.Height()
	if n == 0 {
		return ""
	}
	return m.asciiArt.NextAndString(int(percentageDifference(a, b)))
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var view string
	switch m.state {
	case inputView:
		view = fmt.Sprintf(
			"%s \n\n%s",
			titleStyle.Render(),
			paddingleft.Render(
				fmt.Sprintf("%s\n\n%s", heightThing.Render(fmt.Sprintf("%s \n\n%s", listTitleStyle.Render("Time in minute: "),
					m.input.View(),
				)),
					subtleStyle.Render("↩ Enter")+dotStyle+subtleStyle.Render("Ctrl + C"))))
	case listView:
		view = fmt.Sprintf(
			"%s \n\n%s",
			titleStyle.Render(),
			// listTitleStyle.Render("Select a visual option: "),
			m.list.View(),
		)
	case timerView:
		view = fmt.Sprintf(
			"%s \n\n %s",
			titleStyle.Render(),
			paddingleft.Render(fmt.Sprintf("%s\n%s%s", center(fmt.Sprintf("%d : %d", int(m.timer.Timeout.Minutes()), int(m.timer.Timeout.Seconds())%60), app_width-10),
				m.DrawAscii(m.minute, m.timer.Timeout),
				// m.generated_thing,
				m.helpView())))

	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, appStyle.Render(view))
}

func center(s string, total_width int) string {
	if len(s) >= total_width {
		return s
	}
	n := total_width - len(s)
	div := n / 2
	return strings.Repeat(" ", div) + s
}

func percentageDifference(a, b time.Duration) float64 {
	if a == 0 && b == 0 {
		return 0.0
	}

	return ((a.Seconds() - b.Seconds()) / a.Seconds()) * 100
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if model, ok := m.(model); ok && model.quitting {
		if model.err != nil {
			fmt.Printf("Error Occoured: %s \n", model.err.Error())
			return
		}
		fmt.Printf("\nMinute: %f, Visual Option: %s\n", model.minute.Minutes(), model.selectedItem)
	}
}
