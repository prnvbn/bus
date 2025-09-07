package ui

import (
	"fmt"
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
	etas  []string
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
		mins := int(r.ETA.Minutes())
		if mins < 1 {
			groups[key].etas = append(groups[key].etas, "due")
		} else {
			groups[key].etas = append(groups[key].etas, fmt.Sprintf("%d", mins))
		}
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
		etas := strings.Join(g.etas, ", ")
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
