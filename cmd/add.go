package cmd

import (
	"fmt"
	"iter"
	"os"

	"github.com/prnvbn/bq/internal/tfl"
	"github.com/prnvbn/bq/internal/ui"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	route string
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

		stopPoint := runStopSelection(stopPoints)

		fmt.Printf("✅ Using route: %s\n", route)
		fmt.Printf("✅ Using stop: %s (%s)\n", stopPoint.name, stopPoint.letter)
		fmt.Printf("✅ Using id: %s\n", stopPoint.id)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&route, "route", "r", "", "route name e.g. 135")
}

func runRouteInput() string {
	m := ui.NewInputModal()
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return finalModel.(ui.TextInputModal).Value()
}

type stopItem struct {
	id     string
	name   string
	letter string
}

func (s stopItem) Title() string       { return s.name }
func (s stopItem) Description() string { return "Stop " + s.letter }
func (s stopItem) FilterValue() string { return s.name }

func runStopSelection(stopPoints iter.Seq[tfl.StopPoint]) stopItem {
	stopItemsIter := func(yield func(stopItem) bool) {
		for stopPoint := range stopPoints {
			if !yield(stopItem{name: stopPoint.CommonName, letter: stopPoint.Letter, id: stopPoint.ID}) {
				return
			}
		}
	}

	m := ui.NewListInputModal(stopItemsIter)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	sm := finalModel.(ui.ListInputModel[stopItem])
	return sm.Choice()
}
