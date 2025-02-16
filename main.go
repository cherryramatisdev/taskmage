package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cherryramatisdev/taskmage/cmd"
	"github.com/cherryramatisdev/taskmage/taskwarrior"
	"golang.org/x/term"
)

func main() {
	db, err := taskwarrior.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tasks, err := taskwarrior.GetTasksByStatus(db, taskwarrior.Pending)

	if err != nil {
		panic(err)
	}

	{
		target := time.Now()

		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			panic(err)
		}

		view := cmd.View{
			Width: width,
		}

		fmt.Printf("%s %d %s %d\n", target.Weekday().String(), target.Day(), target.Month().String(), target.Year())
		fmt.Print(view.MountAgendaView(tasks, target))
	}

	fmt.Println() 

	{
		target := time.Now().Add(time.Hour * 24) 

		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			panic(err)
		}

		view := cmd.View{
			Width: width,
		}

		fmt.Printf("%s %d %s %d\n", target.Weekday().String(), target.Day(), target.Month().String(), target.Year())
		fmt.Print(view.MountAgendaView(tasks, target))
	}
}
