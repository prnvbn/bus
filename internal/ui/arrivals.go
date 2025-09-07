package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type ArrivalRow struct {
	Route string
	Stop  string
	ETA   time.Duration
}

type group struct {
	route string
	stop  string
	etas  []float64
}

func RenderArrivals(rows []ArrivalRow) string {
	if len(rows) == 0 {
		return "No arrivals found."
	}

	groups := map[string]*group{}
	for _, r := range rows {
		key := r.Route + "|" + r.Stop
		if _, ok := groups[key]; !ok {
			groups[key] = &group{route: r.Route, stop: r.Stop}
		}
		mins := r.ETA.Minutes()
		groups[key].etas = append(groups[key].etas, mins)
	}

	tableStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("0")).
		Foreground(lipgloss.Color("#FF5F1F")).
		Bold(true).
		PaddingLeft(4).
		PaddingTop(1)

	colRoute := lipgloss.NewStyle().Width(6)
	colStop := lipgloss.NewStyle().Width(33)
	colETA := lipgloss.NewStyle().Width(20)

	var table string
	table += lipgloss.JoinHorizontal(
		lipgloss.Left,
		colRoute.Render("Route"),
		colStop.Render("Stop"),
		colETA.Render("Mins"),
	) + "\n"

	for _, g := range groups {
		sort.Float64s(g.etas)
		parts := make([]string, 0, len(g.etas))
		for _, n := range g.etas {
			if n < 1 {
				parts = append(parts, "due")
			} else {
				parts = append(parts, fmt.Sprintf("%d", int(n)))
			}
		}

		etas := strings.Join(parts, ", ")
		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			colRoute.Render(g.route),
			colStop.Render(g.stop),
			colETA.Render(etas),
		)
		table += row + "\n"
	}

	return tableStyle.Render(table)
}
