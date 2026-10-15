package ui

import (
	"fmt"
	"sort"
	"strconv"
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
	ordered := make([]*group, 0, len(groups))
	for _, g := range groups {
		ordered = append(ordered, g)
	}
	sort.Slice(ordered, func(i, j int) bool {
		a, b := ordered[i], ordered[j]
		if a.route != b.route {
			an, aok := routeNumber(a.route)
			bn, bok := routeNumber(b.route)
			switch {
			case aok && bok && an != bn:
				return an < bn
			case aok && !bok:
				return true
			case !aok && bok:
				return false
			}
			return a.route < b.route
		}
		return a.stop < b.stop
	})

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

	for _, g := range ordered {
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

// routeNumber extracts a number for sorting, trying the full route before
// stripping up to two leading ASCII letters, so "15", "D7", and "SL2" yield
// 15, 7, and 2. It returns 0, false if none of these attempts can be parsed as an int.
func routeNumber(route string) (int, bool) {
	for stripped := 0; stripped <= 2; stripped++ {
		if n, err := strconv.Atoi(route); err == nil {
			return n, true
		}
		if len(route) == 0 || !((route[0] >= 'A' && route[0] <= 'Z') || (route[0] >= 'a' && route[0] <= 'z')) {
			break
		}
		route = route[1:]
	}
	return 0, false
}
