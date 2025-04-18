package main

import (
	"aimssh/ascii_generator"
	"aimssh/helper"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/gen2brain/beeep"
)

type sessionState int

const (
	dotChar   = " • "
	app_width = 50
	host      = "localhost"
	port      = "13234"
	logo      = `    _    ___ __  __   ____ ____  _   _ 
   / \  |_ _|  \/  | / ___/ ___|| | | |
  / _ \  | || |\/| | \___ \___ \| |_| |
 / ___ \ | || |  | |  ___) |__) |  _  |
/_/   \_\___|_|  |_| |____/____/|_| |_|
                                              `

	inputView sessionState = iota
	logoView
	listView
	timerView
)

var (
	appStyle          = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder(), true, true, true, true).Width(app_width)
	heightThing       = lipgloss.NewStyle().Height(21)
	paddingleft       = lipgloss.NewStyle().PaddingLeft(2)
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#49beaa")).Bold(true).SetString(helper.Center(`<AIM SSH>`, app_width-4)).AlignHorizontal(lipgloss.Center)
	listTitleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#bfedc1")).PaddingLeft(-10)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#CFF27E"))

	greenColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#bfedc1")).PaddingLeft(2).Faint(true)
	dotStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	// subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	run_as_ssh bool
	gitlink    = greenColor.Render("https://github.com/sairash/aimssh")

	brownColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#967969"))
	end_info   = fmt.Sprintf("\n Thanks for using %s! \n Give a star %s \n Made By     %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#49beaa")).Bold(true).Render("<AIM SSH>"), gitlink, greenColor.Render("https://sairashgautam.com.np/"))
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
	workingon    textinput.Model
	list         list.Model
	minute       time.Duration
	selectedItem string
	timer        timer.Model
	loadingTimer timer.Model
	keymap       keymap
	help         help.Model
	err          error
	width        int
	height       int
	session      string
	asciiArt     ascii_generator.AsciiArt
	quitting     bool
	timedOut     bool
	sessionSsh   ssh.Session
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
	new   key.Binding

	back  key.Binding
	enter key.Binding
	ctrlc key.Binding
}

func initialModel(session interface{}) model {
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
		item("Tree"),
		item("Flow"),
		item("Coffee"),
	}

	l := list.New(items, itemDelegate{}, 30, 11)
	l.Title = "Select a visual option: "
	l.SetShowStatusBar(false)
	l.Styles.Title = listTitleStyle
	// l.SetShowHelp(false)
	l.SetHeight(23)
	// l.SetShowTitle(false)

	m := model{
		state:        logoView,
		input:        ti,
		loadingTimer: timer.NewWithInterval(800*time.Millisecond, time.Millisecond),
		workingon:    woI,
		list:         l,
		err:          nil,
		session:      "Work",
		timedOut:     false,
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys(" ", "s"),
				key.WithHelp("space", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys(" ", "s"),
				key.WithHelp("space", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "restart"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
			new: key.NewBinding(
				key.WithKeys("n"),
				key.WithHelp("n", "new"),
			),
			back: key.NewBinding(
				key.WithKeys("ctrl+b"),
				key.WithHelp("ctrl+b", "back"),
			),
			enter: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("↩", "enter"),
			),
			ctrlc: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
		},
		help: help.New(),
	}

	switch session := session.(type) {
	case ssh.Session:
		m.sessionSsh = session
	}

	return m

}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.loadingTimer.Init(), m.loadingTimer.Start())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyCtrlB:

			switch m.state {
			case inputView:
				if m.workingon.Focused() {
					m.workingon.Blur()
					m.input.Focus()
				}

			case listView:
				m.state = inputView
				m.input.Blur()
				m.workingon.Focus()
			case timerView:
				m.state = listView
			}

		}

	}

	switch m.state {
	case logoView:
		return updateLogo(msg, m)
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

func updateLogo(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.loadingTimer, cmd = m.loadingTimer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.loadingTimer.Stop()
		m.state = inputView
		return m, textinput.Blink
	}
	return m, nil
}

func updateInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.input.Focused() {
				if m.input.Value() != "" {

					min, err := strconv.Atoi(m.input.Value())

					if err != nil {
						m.err = err
						m.quitting = true
						return m, tea.Quit
					}

					m.minute = time.Duration(min) * time.Minute
					m.input.Blur()
					return m, m.workingon.Focus()
				}
			} else {
				if m.workingon.Value() != "" {
					m.session = m.workingon.Value()
				}
				m.state = listView
			}
		}
	}

	if m.input.Focused() {
		m.input, cmd = m.input.Update(msg)
	} else {
		m.workingon, cmd = m.workingon.Update(msg)
	}

	return m, cmd
}

func (m model) generate_ascii() ascii_generator.AsciiArt {
	switch m.selectedItem {
	case "Coffee":
		return ascii_generator.GenerateCoffee()
	case "Tree":
		return ascii_generator.GenerateTree(40, 18)
	default:
		return ascii_generator.GenerateRow(40, 17)
	}
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
				m.asciiArt = m.generate_ascii()
				// +brownColor.Render(strings.Repeat("░", app_width-8))
				m.keymap.start.SetEnabled(false)

				notificaiton_body := fmt.Sprintf("Timer set for %d minutes.", int(m.timer.Timeout.Minutes()))
				if !run_as_ssh {
					beeep.Alert(helper.TimerStartedTitle, notificaiton_body, "assets/logo.png")
				} else {
					switch m.sessionSsh.Context().Value("operating_system") {
					case "linux":
						wish.Command(m.sessionSsh, "notify-send", "-a", "Aimssh", helper.TimerStartedTitle, notificaiton_body).Run()
					}
				}

				return m, m.timer.Init()
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func format(n int) string {
	var b [2]byte
	if n < 10 {
		b[0], b[1] = '0', byte(n)+'0'
	} else {
		b[0], b[1] = byte(n/10)+'0', byte(n%10)+'0'
	}
	return string(b[:])
}

func updateTimer(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.timedOut = true

		notificaiton_body := "The timer has ended."
		if !run_as_ssh {
			beeep.Alert(helper.TimerEndedTitle, notificaiton_body, "assets/logo.png")
		} else {
			switch m.sessionSsh.Context().Value("operating_system") {
			case "linux":
				wish.Command(m.sessionSsh, "notify-send", "-a", "Aimssh", helper.TimerEndedTitle, notificaiton_body).Run()

			}
		}

		// return m, tea.SetWindowTitle("Done")

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
			m.timedOut = false
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timedOut = false
			m.asciiArt = m.generate_ascii()
			m.timer.Timeout = m.minute

			notificaiton_body := fmt.Sprintf("Timer set for %d minutes.", int(m.timer.Timeout.Minutes()))
			if !run_as_ssh {
				beeep.Alert(helper.TimerRestartedTitle, notificaiton_body, "assets/logo.png")
			} else {
				switch m.sessionSsh.Context().Value("operating_system") {
				case "linux":
					wish.Command(m.sessionSsh, "notify-send", "-a", "Aimssh", helper.TimerRestartedTitle, notificaiton_body).Run()
				}
			}

			return m, m.timer.Start()
		case key.Matches(msg, m.keymap.new):
			m.timedOut = false
			m.state = inputView
			m.workingon.Blur()
			m.input.Focus()
			m.timer.Stop()
			return m, nil
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}

func (m model) helpView() string {
	switch m.state {
	case inputView:
		return "\n" + m.help.ShortHelpView([]key.Binding{
			m.keymap.enter,
			m.keymap.back,
			m.keymap.ctrlc,
		})
	}

	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
		m.keymap.new,
	})
}

func (m model) DrawAscii(a, b time.Duration) string {
	n := m.asciiArt.Height()
	if n == 0 {
		return ""
	}
	if m.selectedItem != "Flow" {
		return m.asciiArt.NextAndString(int(percentageDifference(a, b)))
	}
	return "\n" + greenColor.Render(m.asciiArt.NextAndString(0)) + "\n"
	// camel := `

	// 		  _v   ,,
	// 		 /'|   &&.
	// 		 '-\'-&&&&&.
	// 			 \&&&&&&&\
	// 			 &&#""&& \
	// 			 / /   |\  x
	// 			/ /    / /

	// `
	// return ascii_generator.BrownColor.Render(camel) + "\n"
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var view string
	switch m.state {
	case logoView:
		view = fmt.Sprintf("\n%s  \n\n\n            %s%s\n\n                  Loading...\n", titleStyle.SetString(logo).Render(), greenColor.Bold(true).Render("ssh"), selectedItemStyle.Underline(true).PaddingLeft(2).Render("aim.ftp.sh"))
	case inputView:
		view = fmt.Sprintf(
			"%s \n\n%s",
			titleStyle.Render(),
			paddingleft.Render(
				fmt.Sprintf("%s\n%s",
					heightThing.Render(
						fmt.Sprintf("%s \n\n%s\n\n%s \n\n%s",
							listTitleStyle.Render("Time in minute: *"),
							m.input.View(),
							brownColor.Render("Session:"),
							m.workingon.View(),
						),
					),
					m.helpView())))
	case listView:
		view = fmt.Sprintf(
			"%s \n\n%s%s%s",
			titleStyle.Render(),
			// listTitleStyle.Render("Select a visual option: "),
			m.list.View(),
			dotStyle,
			m.list.Help.ShortHelpView([]key.Binding{m.keymap.back}),
		)
	case timerView:
		text := "Session Ended, Press (r) or (n)"
		if !m.timedOut {
			text = format(int(m.timer.Timeout.Minutes())) + " : " + format(int(m.timer.Timeout.Seconds())%60)
		}
		view = fmt.Sprintf(
			"%s \n\n %s",
			titleStyle.Render(),
			paddingleft.Render(fmt.Sprintf("%s\n%s%s\n%s", helper.Center(text, app_width-10),
				m.DrawAscii(m.minute, m.timer.Timeout),
				helper.Center("Session: "+m.session, app_width-8),
				// m.generated_thing,
				m.helpView())))

	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, appStyle.Render(view))
}

func percentageDifference(a, b time.Duration) float64 {
	if a == 0 && b == 0 {
		return 0.0
	}

	return ((a.Seconds() - b.Seconds()) / a.Seconds()) * 100
}

func wish_server() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
			// After Bubbletea quit.
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					// pty, _, _ :=
					wish.Println(sess, end_info)
					next(sess)
				}
			},
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	s.Context().SetValue("operating_system", detectOS(s))
	return initialModel(s), []tea.ProgramOption{tea.WithAltScreen()}
}

func detectOS(s ssh.Session) string {
	// if wish.Command(s, "sw_vers").Run() == nil {
	// 	return "darwin"
	// }

	if wish.Command(s, "paplay", "--help").Run() == nil {
		return "linux"
	}

	return "windows"
}

func main() {

	flag.BoolVar(&run_as_ssh, "ssh", false, "run as ssh server?")
	flag.Parse()

	if run_as_ssh {
		wish_server()
		return
	}

	p := tea.NewProgram(initialModel(nil), tea.WithAltScreen())
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
		fmt.Println(end_info)
	}
}
