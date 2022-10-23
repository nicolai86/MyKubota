package main

import (
	"log"
	"os"

	"github.com/nicolai86/mykubota/cmd/cli/internal/cmds"

	"github.com/mitchellh/cli"
)

func createCLI(args []string) *cli.CLI {
	c := cli.NewCLI("mykubota", "0.1.0")
	c.Args = args
	c.Commands = map[string]cli.CommandFactory{
		"login": func() (cli.Command, error) {
			return &cmds.LoginCommand{}, nil
		},
		"equipment list": func() (cli.Command, error) {
			return &cmds.EquipmentListCommand{}, nil
		},
		"equipment show": func() (cli.Command, error) {
			return &cmds.EquipmentShowCommand{}, nil
		},
		"equipment maintenance": func() (cli.Command, error) {
			return &cmds.EquipmentMaintenanceCommand{}, nil
		},
	}
	return c
}

func main() {
	c := createCLI(os.Args[1:])

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
