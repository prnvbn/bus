package bus

import (
	"fmt"
	"strings"
	"unicode"
)

func (c Config) ResolveSection(name string) (string, error) {
	name = strings.TrimSpace(name)
	for _, section := range c.Sections {
		if strings.EqualFold(section, name) {
			return section, nil
		}
	}
	return "", fmt.Errorf("section %q does not exist; create it with bus section add %q", name, name)
}

func (c *Config) AddSection(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || strings.ContainsFunc(name, unicode.IsControl) {
		return "", fmt.Errorf("section name must contain text and no control characters")
	}
	if _, err := c.ResolveSection(name); err == nil {
		return "", fmt.Errorf("section %q already exists", name)
	}
	c.Sections = append(c.Sections, name)
	return name, nil
}

func (c *Config) RenameSection(oldName, newName string) error {
	oldName, err := c.ResolveSection(oldName)
	if err != nil {
		return err
	}
	newName = strings.TrimSpace(newName)
	if newName == "" || strings.ContainsFunc(newName, unicode.IsControl) {
		return fmt.Errorf("section name must contain text and no control characters")
	}
	if existing, err := c.ResolveSection(newName); err == nil && existing != oldName {
		return fmt.Errorf("section %q already exists", newName)
	}
	for i, section := range c.Sections {
		if section == oldName {
			c.Sections[i] = newName
		}
	}
	for i := range c.Arrivals {
		if c.Arrivals[i].Section == oldName {
			c.Arrivals[i].Section = newName
		}
	}
	return nil
}

func (c *Config) RemoveSection(name string) error {
	name, err := c.ResolveSection(name)
	if err != nil {
		return err
	}
	for i, section := range c.Sections {
		if section == name {
			c.Sections = append(c.Sections[:i], c.Sections[i+1:]...)
			break
		}
	}
	// A stop may already be unassigned. Keep one copy when removing its section.
	arrivals := make([]Arrival, 0, len(c.Arrivals))
	for _, arrival := range c.Arrivals {
		if arrival.Section == name {
			arrival.Section = ""
		}
		if arrival.Section != "" || !containsArrival(arrivals, arrival) {
			arrivals = append(arrivals, arrival)
		}
	}
	c.Arrivals = arrivals
	return nil
}

func containsArrival(arrivals []Arrival, candidate Arrival) bool {
	for _, arrival := range arrivals {
		if arrival.Route == candidate.Route && arrival.TflID == candidate.TflID && arrival.Section == candidate.Section {
			return true
		}
	}
	return false
}

func (c *Config) AddArrival(arrival Arrival) error {
	if containsArrival(c.Arrivals, arrival) {
		return fmt.Errorf("route %s at %s is already saved in this section", arrival.Route, arrival.StopPoint)
	}
	c.Arrivals = append(c.Arrivals, arrival)
	return nil
}

func (c *Config) MoveArrival(index int, section string) error {
	if index < 0 || index >= len(c.Arrivals) {
		return fmt.Errorf("saved stop does not exist")
	}
	arrival := c.Arrivals[index]
	if arrival.Section == section {
		return nil
	}
	arrival.Section = section
	if containsArrival(c.Arrivals, arrival) {
		return fmt.Errorf("route %s at %s is already saved in the destination section", arrival.Route, arrival.StopPoint)
	}
	c.Arrivals[index] = arrival
	return nil
}

// ValidateSections also normalizes references when a config was edited by hand.
func (c *Config) ValidateSections() error {
	sections := Config{}
	for _, name := range c.Sections {
		if _, err := sections.AddSection(name); err != nil {
			return err
		}
	}
	c.Sections = sections.Sections
	for i := range c.Arrivals {
		if c.Arrivals[i].Section == "" {
			continue
		}
		name, err := c.ResolveSection(c.Arrivals[i].Section)
		if err != nil {
			return err
		}
		c.Arrivals[i].Section = name
	}
	return nil
}
