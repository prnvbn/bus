package cmd

import (
	"fmt"
	"strings"

	"github.com/prnvbn/bus/internal/bus"
	"github.com/prnvbn/bus/internal/tfl"
	"github.com/prnvbn/bus/internal/ui"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a stop to favourites",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		route, err := cmd.Flags().GetString("route")
		if err != nil {
			return err
		}
		section, err := cmd.Flags().GetString("section")
		if err != nil {
			return err
		}
		sectionSet := cmd.Flags().Changed("section")
		if sectionSet && strings.TrimSpace(section) != "" {
			section, err = cfg.ResolveSection(section)
			if err != nil {
				return err
			}
		} else {
			section = ""
		}

		if route == "" {
			route, err = ui.NewInputModal("🚌 Enter Bus Route", "for e.g. 135 or D7").Run()
			if err != nil {
				return sectionSelectionError(cmd, err)
			}
		}
		route = strings.ToUpper(strings.TrimSpace(route))
		if route == "" {
			return fmt.Errorf("bus route cannot be empty")
		}

		c := tfl.NewClient()
		stopPoints, err := c.StopPoints(route)
		if err != nil {
			return fmt.Errorf("get stop points: %w", err)
		}
		stopItemsIter := func(yield func(stopItem) bool) {
			for stopPoint := range stopPoints {
				if !yield(stopItem{name: stopPoint.CommonName, letter: stopPoint.Letter, id: stopPoint.ID}) {
					return
				}
			}
		}
		stopPoint, err := ui.NewListInputModal(stopItemsIter, "Choose a stop").Run()
		if err != nil {
			return sectionSelectionError(cmd, err)
		}
		if !sectionSet && len(cfg.Sections) > 0 {
			section, err = selectSection("Add this stop to:")
			if err != nil {
				return sectionSelectionError(cmd, err)
			}
		}
		if err := cfg.AddArrival(bus.Arrival{
			Route: route, StopPoint: stopPoint.name, Letter: stopPoint.letter,
			TflID: stopPoint.id, Section: section,
		}); err != nil {
			return err
		}
		cmd.Printf("Saved route %s at %s to %s.\n", route, stopPoint.name, sectionLabel(section))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("route", "r", "", "route name e.g. 135")
	addCmd.Flags().String("section", "", "save to an existing section; use an empty value for no section")
}

type stopItem struct {
	id     string
	name   string
	letter string
}

func (s stopItem) Title() string       { return s.name }
func (s stopItem) Description() string { return "Stop " + s.letter }
func (s stopItem) FilterValue() string { return s.name }
