package cmd

import (
	"fmt"
	"iter"
	"os"

	"github.com/prnvbn/bq/internal/bq"
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

		var err error
		if route == "" {
			route, err = ui.NewInputModal(
				"🚌 Enter Bus Route",
				"for e.g. 135 or D7",
			).Run()
			fatal(err, "error running route input")
		}

		stopPoints, err := c.StopPoints(route)
		fatal(err, "error getting stop points")

		stopPoint := runStopSelection(stopPoints)

		cfg.Arrivals = append(cfg.Arrivals, bq.Arrival{
			Route:     route,
			StopPoint: stopPoint.name,
			Letter:    stopPoint.letter,
			TflID:     stopPoint.id,
		})
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&route, "route", "r", "", "route name e.g. 135")
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
