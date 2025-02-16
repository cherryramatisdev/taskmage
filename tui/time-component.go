package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type level int

const (
	hourLevel level = iota
	minuteLevel
	secondLevel
)

var levels = []level{hourLevel, minuteLevel, secondLevel}

func (l *level) next() {
	curr := int(*l)
	*l = levels[(curr+1)%len(levels)]
}

func (l *level) prev() {
	curr := int(*l)
	*l = levels[(curr+len(levels)-1)%len(levels)]
}

type hourStyles struct {
	borderStyle   lipgloss.Style
	focusedStyle  lipgloss.Style
	fullTextStyle lipgloss.Style
}

type HourSelectorModel struct {
	Hour     int
	Minute   int
	Second   int
	selected level
	styles   hourStyles
}

func (m *HourSelectorModel) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", m.Hour, m.Minute, m.Second)
}

func (m *HourSelectorModel) inc() {
	switch m.selected {
	case hourLevel:
		m.Hour = (m.Hour + 1) % 24
	case minuteLevel:
		m.Minute = (m.Minute + 1) % 60
	case secondLevel:
		m.Second = (m.Second + 1) % 60
	}
}

func (m *HourSelectorModel) desc() {
	switch m.selected {
	case hourLevel:
		m.Hour = (m.Hour + 23) % 24
	case minuteLevel:
		m.Minute = (m.Minute + 59) % 60
	case secondLevel:
		m.Second = (m.Second + 59) % 60
	}
}

func (m *HourSelectorModel) Init() tea.Cmd {
	return nil
}

func (m *HourSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "enter":
			return m, tea.Quit
		case "up", "k":
			m.inc()
			return m, nil
		case "down", "j":
			m.desc()
			return m, nil
		case "l", "right":
			m.selected.next()
			return m, nil
		case "left", "h":
			m.selected.prev()
			return m, nil
		}
	}
	return m, nil
}

func (m *HourSelectorModel) View() string {
	format := func(v int, highlight bool) string {
		str := fmt.Sprintf("%02d", v)
		if highlight {
			return m.styles.focusedStyle.Render(str)
		}
		return str
	}

	ret := m.styles.borderStyle.Render(m.styles.fullTextStyle.Render(
		fmt.Sprintf("%s:%s:%s",
			format(m.Hour, m.selected == hourLevel),
			format(m.Minute, m.selected == minuteLevel),
			format(m.Second, m.selected == secondLevel),
		)))
	return fmt.Sprintf("%s\n", ret)
}

func NewHourSelectorModel() HourSelectorModel {
	bs := lipgloss.NewStyle().Width(26).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("212"))
	fts := lipgloss.NewStyle().PaddingLeft(9)
	fs := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)

	return HourSelectorModel{
		Hour:     0,
		Minute:   0,
		Second:   0,
		selected: hourLevel,
		styles: hourStyles{
			focusedStyle:  fs,
			fullTextStyle: fts,
			borderStyle:   bs,
		},
	}
}
