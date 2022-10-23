package cmds

import (
	"context"
	"flag"
	"log"

	"github.com/nicolai86/mykubota/cmd/cli/internal/auth"

	"github.com/mitchellh/cli"
	"github.com/nicolai86/mykubota"
	"github.com/posener/complete"
)

type LoginCommand struct {
}

var _ cli.Command = &LoginCommand{}

func (c *LoginCommand) Help() string {
	return `Login performs authentication with the MyKubota API. It is required to access private information like the list of equipment.
required flags:
	-u, --username MyKubota username
	-p, --password MyKubota password`
}

func (c *LoginCommand) Run(args []string) int {
	var username, password string
	flag.StringVar(&username, "username", "", "mykubota username")
	flag.StringVar(&username, "u", "", "mykubota username")
	flag.StringVar(&password, "password", "", "mykubota password")
	flag.StringVar(&password, "p", "", "mykubota password")

	if err := flag.CommandLine.Parse(args); err != nil {
		return -1
	}
	if username == "" || password == "" {
		return cli.RunResultHelp
	}
	client := mykubota.New("en-CA")
	session, err := client.Authenticate(context.Background(), username, password)
	if err != nil {
		log.Fatalf("Error authenticating: %s", err)
	}

	if err := auth.Persist(session); err != nil {
		log.Fatalf("Error persisting session: %s", err)
	}

	return 0
}

func (c *LoginCommand) Synopsis() string {
	return "perform authentication with the MyKubota API"
}

func (c *LoginCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"--username": complete.PredictAnything,
		"--password": complete.PredictAnything,
		"-u":         complete.PredictAnything,
		"-p":         complete.PredictAnything,
	}
}
