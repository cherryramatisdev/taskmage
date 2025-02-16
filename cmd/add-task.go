package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/cherryramatisdev/taskmage/taskwarrior"
	"github.com/cherryramatisdev/taskmage/tui"
	"github.com/spf13/cobra"
)

type view int

const (
	formView view = iota
	datetimeView
)

type model struct {
	description string
	dueDate     time.Time
	tags        []string

	datetimeInput tui.DateAndHourModel
	form          *huh.Form

	curView  view
	quitting bool
}

func initialModel(tags []string) model {
	form := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Key("description").
			Title("Description"),
		huh.NewMultiSelect[string]().
			Key("tags").
			Title("Tags").
			Options(huh.NewOptions(tags...)...),
	))

	return model{
		datetimeInput: tui.NewDateAndHourModel(),
		form:          form,
		curView:       formView,
	}
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.curView {
	case formView:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}

		if m.form.State == huh.StateCompleted {
			m.curView = datetimeView
		}

		return m, tea.Batch(cmds...)
	case datetimeView:
		m.datetimeInput, cmd = m.datetimeInput.Update(msg)

		if m.datetimeInput.Finished {
			m.dueDate = m.datetimeInput.Time()
			m.description = m.form.GetString("description")
			m.tags = m.form.Get("tags").([]string)

			m.quitting = true

			var tags strings.Builder

			for _, v := range m.tags {
				tags.WriteString(fmt.Sprintf("+%s ", v))
			}

			cmd := exec.Command("task", "add", m.description, fmt.Sprintf("due:'%d-%02d-%02dT%02d:%02d'", m.dueDate.Year(), m.dueDate.Month(), m.dueDate.Day(), m.dueDate.Hour(), m.dueDate.Minute()), tags.String())

			stdout, err := cmd.Output()

			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Print(string(stdout))

			return m, tea.Quit
		}

		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	if m.quitting {
		s += "\n"
	}

	switch m.curView {
	case formView:
		s = m.form.View()
	case datetimeView:
		s = m.datetimeInput.View()
	default:
		s = ""
	}

	return s
}

func RegisterAddTaskCmd(db *sql.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "This command will open a tui like experience for you to inform the title of the task, tags and due date with an interactive calendar",
		Run: func(cmd *cobra.Command, args []string) {
			tags, err := taskwarrior.GetTags(db)

			if err != nil {
				fmt.Fprintf(os.Stderr, "deu ruim no db: %v", err)
				os.Exit(1)
			}

			p := tea.NewProgram(initialModel(tags))
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		},
	}
}
