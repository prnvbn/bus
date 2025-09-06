package cmd

import (
	"errors"
	"fmt"

	"github.com/prnvbn/bq/internal/bq"
	"github.com/prnvbn/bq/internal/tfl"
	"github.com/prnvbn/bq/internal/ui"
	"github.com/spf13/cobra"
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

			if errors.Is(err, ui.ErrCancelled) {
				fmt.Println("User cancelled route input")
				return
			}
			fatal(err, "error running route input")
		}

		stopPoints, err := c.StopPoints(route)
		fatal(err, "error getting stop points")

		stopItemsIter := func(yield func(stopItem) bool) {
			for stopPoint := range stopPoints {
				if !yield(stopItem{name: stopPoint.CommonName, letter: stopPoint.Letter, id: stopPoint.ID}) {
					return
				}
			}
		}
		stopPoint, err := ui.NewListInputModal(stopItemsIter).Run()
		if errors.Is(err, ui.ErrCancelled) {
			fmt.Println("User cancelled stop selection")
			return
		}
		fatal(err, "error running stop selection")

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
