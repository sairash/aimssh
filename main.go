package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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
	inputView sessionState = iota
	listView
)

var (
	appStyle          = lipgloss.NewStyle().Padding(1, 2)
	heightThing       = lipgloss.NewStyle().Height(9)
	paddingleft       = lipgloss.NewStyle().PaddingLeft(2)
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#49beaa")).Bold(true).SetString("Zen Cli \n\n")
	listTitleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#bfedc1")).PaddingLeft(-10)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#CFF27E"))
	dotStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	subtleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
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
	minute       int
	selectedItem string
	timer        timer.Model
	err          error
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
	// l.SetShowTitle(false)

	return model{
		state: inputView,
		input: ti,
		list:  l,
		err:   nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

				m.minute = min
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
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
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
	}

	return appStyle.Render(view)
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
		fmt.Printf("\nMinute: %d, Visual Option: %s\n", model.minute, model.selectedItem)
	}
}
