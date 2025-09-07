package cmd

import (
	"fmt"
	"time"

	"github.com/prnvbn/bq/internal/tfl"
	"github.com/prnvbn/bq/internal/ui"
	"github.com/spf13/cobra"
)

var arrivalsCmd = &cobra.Command{
	Use: "arrivals",
	Run: func(cmd *cobra.Command, args []string) {
		c := tfl.NewClient()

		rows := []ui.ArrivalRow{}

		for _, configured := range cfg.Arrivals {
			arrivals, err := c.Arrivals(configured.TflID)
			fatal(err, "error getting arrivals %w", err)

			for a := range arrivals {
				if a.LineName != configured.Route {
					continue
				}
				rows = append(rows, ui.ArrivalRow{
					Route: a.LineName,
					Stop:  fmt.Sprintf("%s (%s)", a.StationName, a.PlatformName),
					ETA:   time.Duration(a.TimeToStation) * time.Second,
				})
			}
		}

		fmt.Println(ui.RenderArrivals(rows))
	},
}

func init() {
	rootCmd.AddCommand(arrivalsCmd)
}
