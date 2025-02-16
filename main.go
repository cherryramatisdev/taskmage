package main

import (
	"fmt"
	"os"

	"github.com/cherryramatisdev/taskmage/cmd"
	"github.com/cherryramatisdev/taskmage/taskwarrior"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "taskmage",
	Short: "taskmage is a suite around taskwarrior to help you better manage your tasks",
}

func main() {
	db, err := taskwarrior.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cmd.RegisterAgendaCmd(&cmd.View{Width: width, DB: db}))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 	tasks, err := taskwarrior.GetTasksByStatus(db, taskwarrior.Pending)

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	{
	// 		target := time.Now()

	// 		view := cmd.View{
	// 			Width: width,
	// 		}

	// 		fmt.Printf("%s %d %s %d\n", target.Weekday().String(), target.Day(), target.Month().String(), target.Year())
	// 		fmt.Print(view.MountAgendaView(tasks, target))
	// 	}

	// 	fmt.Println()

	// 	{
	// 		target := time.Now().Add(time.Hour * 24)

	// 		width, _, err := term.GetSize(int(os.Stdout.Fd()))
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		view := cmd.View{
	// 			Width: width,
	// 		}

	//		fmt.Print(view.MountAgendaView(tasks, target))
	//	}
}
