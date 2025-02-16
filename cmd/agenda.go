package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cherryramatisdev/taskmage/taskwarrior"
	"github.com/cherryramatisdev/taskmage/tui"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const INITIAL_HOUR int = 8
const FINAL_HOUR int = 23

// Week-agenda (W07):
// Monday     10 February 2025 W07
// Tuesday    11 February 2025
// Wednesday  12 February 2025
// Thursday   13 February 2025
// Friday     14 February 2025
// Saturday   15 February 2025
//                8:00 ┄┄┄┄┄ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
//   refile:      8:30 ┄┄┄┄┄ Scheduled:  TODO teste teste
//               10:00 ┄┄┄┄┄ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
//               12:00 ┄┄┄┄┄ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
//               14:00 ┄┄┄┄┄ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
//               16:00 ┄┄┄┄┄ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
//               18:00 ┄┄┄┄┄ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
//               20:00 ┄┄┄┄┄ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
//               21:32 ┄┄┄┄┄ ← now ───────────────────────────────────────────────
//   refile:     22:32 ┄┄┄┄┄ Scheduled:  TODO teste testando
//   refile:     Sched.51x:  TODO work on the article
// Sunday     16 February 2025

type View struct {
	Width int
	DB    *sql.DB
}

func (v *View) MountAgendaView(tasks []*taskwarrior.Task, target time.Time) string {
	var output strings.Builder

	line := tui.DrawLine(v.Width / 2)

	for i := INITIAL_HOUR; i <= FINAL_HOUR; i++ {
		var hour int

		if i > 23 {
			hour = 0
		} else {
			hour = i
		}

		filteredTasks := taskwarrior.FindTaskByDueDate(tasks, time.Date(target.Year(), target.Month(), target.Day(), hour, 0, 0, 0, time.UTC))

		output.WriteString(fmt.Sprintf("%02d:00 %s\n", i, line))

		if target.Hour() == i {
			output.WriteString(color.CyanString(fmt.Sprintf("%02d:%02d %s ← now\n", target.Hour(), target.Minute(), line)))
		}

		if filteredTasks != nil {
			for _, task := range filteredTasks {
				output.WriteString(color.GreenString(fmt.Sprintf("%02d:%02d %s\n", task.Due.Hour(), task.Due.Minute(), task.Description)))
			}
		}
	}

	return output.String()
}

func getDateByWeekday(now time.Time, weekday time.Weekday) time.Time {
	daysUntilTarget := (int(now.Weekday()) - int(weekday) + 7) % 7
	return now.AddDate(0, 0, -daysUntilTarget)
}

func fromStringToWeekday(day string) time.Weekday {
	switch day {
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	case "sunday":
		return time.Sunday
	default:
		return time.Monday
	}
}

func fromHumanDayIdentifierToTimeTarget(day string) time.Time {
	now := time.Now()
	if day == "today" {
		return now
	}

	if day == "tomorrow" {
		return now.Add(time.Hour * 24)
	}

	return getDateByWeekday(now, fromStringToWeekday(day))
}

func RegisterAgendaCmd(view *View) *cobra.Command {
	return &cobra.Command{
		Use:   "agenda [<today|tomorrow|monday|tuesday|wednesday|thursday|friday|saturday|sunday>]",
		Short: "Output a daily resume for the tasks, without any flag will display agenda for current day",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := taskwarrior.GetTasksByStatus(view.DB, taskwarrior.Pending)

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			target := fromHumanDayIdentifierToTimeTarget(args[0])

			fmt.Printf("%s %d %s %d\n", target.Weekday().String(), target.Day(), target.Month().String(), target.Year())
			fmt.Print(view.MountAgendaView(tasks, target))
		},
	}
}
