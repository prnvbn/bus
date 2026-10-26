package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/prnvbn/bus/internal/tfl"
	"github.com/prnvbn/bus/internal/ui"
	"github.com/spf13/cobra"
)

var arrivalsCmd = &cobra.Command{
	Use:   "arrivals",
	Short: "show the arrival times of the tracked buses",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		section, err := cmd.Flags().GetString("section")
		if err != nil {
			return err
		}
		filtered := cmd.Flags().Changed("section")
		if filtered && strings.TrimSpace(section) != "" {
			section, err = cfg.ResolveSection(section)
			if err != nil {
				return err
			}
		} else {
			section = ""
		}

		c := tfl.NewClient()
		rows := map[string][]ui.ArrivalRow{}
		// One request per stop, even when it appears in several sections or routes.
		byStop := map[string][]tfl.Arrival{}
		for _, configured := range cfg.Arrivals {
			if filtered && configured.Section != section {
				continue
			}
			arrivals, fetched := byStop[configured.TflID]
			if !fetched {
				results, err := c.Arrivals(configured.TflID)
				if err != nil {
					return fmt.Errorf("get arrivals for %s: %w", configured.StopPoint, err)
				}
				for arrival := range results {
					arrivals = append(arrivals, arrival)
				}
				byStop[configured.TflID] = arrivals
			}
			for _, a := range arrivals {
				if a.LineName != configured.Route {
					continue
				}
				rows[configured.Section] = append(rows[configured.Section], ui.ArrivalRow{
					Route: a.LineName,
					Stop:  fmt.Sprintf("%s (%s)", a.StationName, a.PlatformName),
					ETA:   time.Duration(a.TimeToStation) * time.Second,
				})
			}
		}
		if len(cfg.Sections) == 0 {
			cmd.Println(ui.RenderArrivals(rows[""]))
			return nil
		}
		sections := append([]string(nil), cfg.Sections...)
		for _, arrival := range cfg.Arrivals {
			if arrival.Section == "" {
				sections = append(sections, "")
				break
			}
		}
		if filtered {
			sections = []string{section}
		}
		groups := make([]ui.ArrivalSection, 0, len(sections))
		for _, name := range sections {
			groups = append(groups, ui.ArrivalSection{Name: sectionLabel(name), Rows: rows[name]})
		}
		cmd.Println(ui.RenderArrivalSections(groups))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(arrivalsCmd)
	arrivalsCmd.Flags().String("section", "", "show only this section; use an empty value for unassigned stops")
	arrivalsCmd.Flags().BoolP("keep-open", "k", false, "keep the arrivals display open and updating (every 5 seconds)")
}
