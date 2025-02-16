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

	rootCmd.AddCommand(cmd.RegisterAgendaCmd(&cmd.View{Width: width, DB: db}), cmd.RegisterAddTaskCmd(db))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
