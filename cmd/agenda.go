package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/cherryramatisdev/taskmage/taskwarrior"
	"github.com/cherryramatisdev/taskmage/tui"
	"github.com/fatih/color"
)

const INITIAL_HOUR int = 8

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
}

func (v *View) MountAgendaView(tasks []*taskwarrior.Task, target time.Time) string {
	var output strings.Builder


	line := tui.DrawLine(v.Width / 2)

	var finalHour int

	if target.Hour() <= INITIAL_HOUR {
		finalHour = 24
	} else {
		finalHour = target.Hour()
	}

	for i := INITIAL_HOUR; i <= finalHour; i++ {
  	var hour int

  	if i > 23 {
    	hour = 0
  	} else {
    	hour = i
  	}

		filteredTasks := taskwarrior.FindTaskByDueDate(tasks, time.Date(target.Year(), target.Month(), target.Day(), hour, 0, 0, 0, time.UTC))

		output.WriteString(fmt.Sprintf("%02d:00 %s\n", i, line))

		if filteredTasks != nil {
			for _, task := range filteredTasks {
				output.WriteString(color.GreenString(fmt.Sprintf("%02d:%02d %s\n", task.Due.Hour(), task.Due.Minute(), task.Description)))
			}
		}
	}

	output.WriteString(fmt.Sprintf("%02d:%02d %s ← now\n", target.Hour(), target.Minute(), line))

	return output.String()
}
