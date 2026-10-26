package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

const (
	arrivalRouteWidth = 6
	arrivalStopWidth  = 33
	arrivalETAWidth   = 20
)

type ArrivalRow struct {
	Route string
	Stop  string
	ETA   time.Duration
}

type ArrivalSection struct {
	Name string
	Rows []ArrivalRow
}

func RenderArrivalSections(sections []ArrivalSection) string {
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#8BD5CA")).Bold(true)
	rule := lipgloss.NewStyle().Foreground(lipgloss.Color("#6C7086"))
	indent := lipgloss.NewStyle().PaddingLeft(4)
	groups := make([]string, 0, len(sections))
	for _, section := range sections {
		width := arrivalRouteWidth + arrivalStopWidth + arrivalETAWidth
		heading := rule.Render("── ") + title.Render(section.Name) +
			rule.Render(" "+strings.Repeat("─", max(2, width-lipgloss.Width(section.Name)-4)))
		body := renderArrivals(section.Rows, 0)
		if len(section.Rows) == 0 {
			body = indent.Render(body)
		}
		groups = append(groups, indent.Render(heading)+"\n"+strings.TrimRight(body, " \t\r\n"))
	}
	return lipgloss.NewStyle().PaddingTop(1).Render(strings.Join(groups, "\n\n"))
}

type group struct {
	route string
	stop  string
	etas  []float64
}

func RenderArrivals(rows []ArrivalRow) string {
	return renderArrivals(rows, 1)
}

func renderArrivals(rows []ArrivalRow, paddingTop int) string {
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
		Foreground(lipgloss.Color("#FF5F1F")).
		Bold(true).
		PaddingLeft(4).
		PaddingTop(paddingTop)

	colRoute := lipgloss.NewStyle().Width(arrivalRouteWidth)
	colStop := lipgloss.NewStyle().Width(arrivalStopWidth)
	colETA := lipgloss.NewStyle().Width(arrivalETAWidth)

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
