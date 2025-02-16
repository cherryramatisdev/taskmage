package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	datepicker "github.com/ethanefung/bubble-datepicker"
)

type DateAndHourModel struct {
	level      dateAndHourLevel
	datepicker datepicker.Model
	hours      HourSelectorModel
	Finished   bool
}

var faint = lipgloss.NewStyle().Faint(true)

func (m DateAndHourModel) Time() time.Time {
	date := m.datepicker.Time
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		m.hours.Hour,
		m.hours.Minute,
		m.hours.Second,
		0,
		time.Local,
	)
}

func (m DateAndHourModel) String() string {
	return m.Time().Format("2006-01-02 15:04:05")
}

func NewDateAndHourModel() DateAndHourModel {
	datepicker := datepicker.New(time.Now())
	datepicker.SelectDate()
	return DateAndHourModel{
		datepicker: datepicker,
		hours:      NewHourSelectorModel(),
		level:      dateLevel,
	}
}

type dateAndHourLevel int

const (
	dateLevel dateAndHourLevel = iota
	hoursLevel
)

var dateAndHourLevels = []dateAndHourLevel{dateLevel, hoursLevel}

func (l *dateAndHourLevel) next() {
	curr := int(*l)
	*l = dateAndHourLevels[(curr+1)%len(dateAndHourLevels)]
}

func (l *dateAndHourLevel) prev() {
	curr := int(*l)
	*l = dateAndHourLevels[(curr+len(dateAndHourLevels)-1)%len(dateAndHourLevels)]
}

func (m DateAndHourModel) Init() tea.Cmd {
	return nil
}

func (m DateAndHourModel) Update(msg tea.Msg) (DateAndHourModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.level == hoursLevel {
  			m.Finished = true
				return m, nil
			}
			m.level.next()
			return m, nil
		case "delete", "backspace":
			if m.level == hoursLevel {
				m.level.prev()
				return m, nil
			}
			return m, nil
		default:
			return m.updateBasedOnLevel(msg)
		}
	}
	return m, nil
}

func (m DateAndHourModel) View() string {
	return fmt.Sprintf("%s\n%s\n", m.datepicker.View(), m.hours.View())
}

func (m DateAndHourModel) updateBasedOnLevel(msg tea.Msg) (DateAndHourModel, tea.Cmd) {
	switch m.level {
	case dateLevel:
		datepicker, cmd := m.datepicker.Update(msg)
		m.datepicker = datepicker
		return m, cmd
	case hoursLevel:
		hours, cmd := m.hours.Update(msg)
		m.hours = *hours.(*HourSelectorModel)
		return m, cmd
	}
	return m, nil
}
