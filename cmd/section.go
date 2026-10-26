package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/prnvbn/bus/internal/bus"
	"github.com/prnvbn/bus/internal/ui"
	"github.com/spf13/cobra"
)

var sectionCmd = &cobra.Command{
	Use: "section", Short: "Organise saved stops into optional sections",
}

func init() {
	rootCmd.AddCommand(sectionCmd)
	sectionCmd.AddCommand(
		&cobra.Command{
			Use: "add <name>", Short: "Create a section", Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				name, err := cfg.AddSection(args[0])
				if err != nil {
					return err
				}
				cmd.Printf("Created section %q.\n", name)
				return nil
			},
		},
		&cobra.Command{
			Use: "list", Short: "List sections and saved stop counts", Args: cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				if len(cfg.Sections) == 0 {
					cmd.Println("No sections. Create one with bus section add Home.")
					return
				}
				for _, section := range cfg.Sections {
					count := 0
					for _, arrival := range cfg.Arrivals {
						if arrival.Section == section {
							count++
						}
					}
					cmd.Printf("%s\t%d saved stops\n", section, count)
				}
			},
		},
		&cobra.Command{
			Use: "rename <name> <new-name>", Short: "Rename a section", Args: cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := cfg.RenameSection(args[0], args[1]); err != nil {
					return err
				}
				cmd.Printf("Renamed section to %q.\n", strings.TrimSpace(args[1]))
				return nil
			},
		},
		&cobra.Command{
			Use: "remove <name>", Short: "Remove a section, keeping its stops unassigned", Args: cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := cfg.RemoveSection(args[0]); err != nil {
					return err
				}
				cmd.Println("Removed section. Its stops are now under Other stops; duplicate stops were merged.")
				return nil
			},
		},
		&cobra.Command{
			Use: "move", Short: "Choose a saved stop and move it to a section", Args: cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(cfg.Arrivals) == 0 {
					cmd.Println("No saved stops. Add one with bus add.")
					return nil
				}
				items := make([]savedStopItem, 0, len(cfg.Arrivals))
				for i, arrival := range cfg.Arrivals {
					items = append(items, savedStopItem{arrival: arrival, index: i})
				}
				stop, err := ui.NewListInputModal(slices.Values(items), "Choose a saved stop to move").Run()
				if err != nil {
					return sectionSelectionError(cmd, err)
				}
				section, err := selectSection("Move this stop to:")
				if err != nil {
					return sectionSelectionError(cmd, err)
				}
				if err := cfg.MoveArrival(stop.index, section); err != nil {
					return err
				}
				cmd.Printf("Moved route %s at %s to %s.\n", stop.arrival.Route, stop.arrival.StopPoint, sectionLabel(section))
				return nil
			},
		},
	)
}

type sectionItem struct {
	name   string
	create bool
}

func (s sectionItem) Title() string {
	if s.create {
		return "+ Create section"
	}
	if s.name == "" {
		return "No section"
	}
	return s.name
}

func (s sectionItem) Description() string {
	if s.create {
		return "Enter a new section name"
	}
	if s.name == "" {
		return "Keep this stop under Other stops"
	}
	return "Saved section"
}

func (s sectionItem) FilterValue() string { return s.Title() }

func selectSection(title string) (string, error) {
	items := make([]sectionItem, 0, len(cfg.Sections)+2)
	for _, section := range cfg.Sections {
		items = append(items, sectionItem{name: section})
	}
	items = append(items, sectionItem{}, sectionItem{create: true})
	choice, err := ui.NewListInputModal(slices.Values(items), title).Run()
	if err != nil {
		return "", err
	}
	if !choice.create {
		return choice.name, nil
	}
	name, err := ui.NewInputModal("Section name", "e.g. Home or Work").Run()
	if err != nil {
		return "", err
	}
	return cfg.AddSection(name)
}

func sectionSelectionError(cmd *cobra.Command, err error) error {
	if errors.Is(err, ui.ErrCancelled) {
		cmd.Println("Cancelled. No changes saved.")
		return nil
	}
	return err
}

func sectionLabel(section string) string {
	if section == "" {
		return "Other stops"
	}
	return section
}

type savedStopItem struct {
	arrival bus.Arrival
	index   int
}

func (s savedStopItem) Title() string {
	return fmt.Sprintf("%s · %s", s.arrival.Route, s.arrival.StopPoint)
}

func (s savedStopItem) Description() string {
	return fmt.Sprintf("Stop %s · %s", s.arrival.Letter, sectionLabel(s.arrival.Section))
}

func (s savedStopItem) FilterValue() string { return s.Title() + " " + s.Description() }
