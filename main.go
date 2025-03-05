package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	inputView sessionState = iota
	listView
)

var (
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
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
	name         string
	selectedItem string
	quitting     bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "10"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	items := []list.Item{
		item("Tree"),
		item("Tomato"),
		item("Coffee"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
		item("Carrot"),
	}

	l := list.New(items, itemDelegate{}, 30, 10)
	l.Title = "Select a visual option: "
	l.SetShowStatusBar(false)
	// l.SetShowTitle(false)

	return model{
		state: inputView,
		input: ti,
		list:  l,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(30, 10)
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
				m.name = m.input.Value()
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
			"Zen cli: \n\nTime in minute: \n\n%s\n\n%s",
			m.input.View(),
			"(Press Enter to continue)",
		)
	case listView:
		view = fmt.Sprintf(
			"Zen cli: \n\n%s\n\n",
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
		fmt.Printf("\nHello %s! You selected: %s\n", model.name, model.selectedItem)
	}
}
