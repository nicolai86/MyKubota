package cmds

import (
	"context"
	"fmt"
	"github.com/mitchellh/cli"
	"log"

	"github.com/nicolai86/mykubota/cmd/cli/internal/auth"
)

type EquipmentListCommand struct{}

var _ cli.Command = &EquipmentListCommand{}

// Help prints equipment list help
func (c *EquipmentListCommand) Help() string {
	return `List equipment.`
}

// Run lists equipment
func (c *EquipmentListCommand) Run(args []string) int {
	sess, err := auth.Restore()
	if err != nil {
		log.Fatalf("Error restoring session: %s", err)
	}

	equipment, err := sess.ListEquipment(context.Background())
	if err != nil {
		log.Fatalf("Error listing equipment: %s", err)
	}

	for _, e := range equipment {
		fmt.Printf("%s (%s, %s, sn %s)\n", e.Nickname, e.Model, e.ID, e.Serial)
	}

	return 0
}

// Synopsis returns the equipment list command synopsis
func (c *EquipmentListCommand) Synopsis() string {
	return "list equipment"
}
