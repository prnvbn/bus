package cmd

import (
	"fmt"
	"iter"
	"os"

	"github.com/prnvbn/bq/internal/tfl"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	route    string
	stopName string
	letter   string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a stop to favourites",
	Run: func(cmd *cobra.Command, args []string) {
		c := tfl.NewClient()

		if route == "" {
			route = runRouteInput()
		}

		stopPoints, err := c.StopPoints(route)
		fatal(err, "error getting stop points")

		if stopName == "" || letter == "" {
			stopPoint := runStopSelection(stopPoints)

			fmt.Printf("🛑 Stop: %s (%s)\n", stopPoint.name, stopPoint.letter)
		}

		fmt.Printf("✅ Using route: %s\n", route)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&route, "route", "r", "", "route name e.g. 135")
	addCmd.Flags().StringVarP(&stopName, "stop", "s", "", "stop point e.g. Primrose Street")
	addCmd.Flags().StringVarP(&letter, "letter", "l", "", "stop letter e.g. J")
}

// -------- Route Input (text) --------
type routeModel struct {
	input textinput.Model
	done  bool
}

func initialRouteModel() routeModel {
	ti := textinput.New()
	ti.Placeholder = "e.g. 135"
	ti.Width = len(ti.Placeholder) + 5
	ti.Focus()

	return routeModel{input: ti}
}

func (m routeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m routeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m routeModel) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf(
		"🚌 Enter Bus Route\n\n%s\n\n(press Enter to confirm)\n",
		m.input.View(),
	)
}

func runRouteInput() string {
	m := initialRouteModel()
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return finalModel.(routeModel).input.Value()
}

// -------- Stop Selection (list) --------
type stopItem struct {
	id     string
	name   string
	letter string
}

func (s stopItem) Title() string       { return s.name }
func (s stopItem) Description() string { return "Stop " + s.letter }
func (s stopItem) FilterValue() string { return s.name }

type stopModel struct {
	list   list.Model
	choice stopItem
}

func newStopModel(stops iter.Seq[tfl.StopPoint]) stopModel {
	items := make([]list.Item, 0)
	for stop := range stops {
		stopItem := stopItem{name: stop.CommonName, letter: stop.Letter, id: stop.ID}
		items = append(items, stopItem)
	}

	l := list.New(items, list.NewDefaultDelegate(), 40, 0)
	l.Title = "Choose a Stop"

	return stopModel{list: l}
}

func (m stopModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m stopModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	preFilterState := m.list.FilterState()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if preFilterState != list.Filtering {
				if item, ok := m.list.SelectedItem().(stopItem); ok {
					m.choice = item
				}
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m stopModel) View() string {
	return m.list.View()
}

func runStopSelection(stopPoints iter.Seq[tfl.StopPoint]) stopItem {
	m := newStopModel(stopPoints)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	sm := finalModel.(stopModel)
	return sm.choice
}
