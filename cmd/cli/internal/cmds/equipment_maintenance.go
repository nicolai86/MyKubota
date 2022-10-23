package cmds

import (
	"context"
	"github.com/mitchellh/cli"
	"github.com/nicolai86/mykubota/cmd/cli/internal/auth"
	"log"
)

type EquipmentMaintenanceCommand struct{}

var _ cli.Command = &EquipmentMaintenanceCommand{}

// Help prints equipment show help
func (c *EquipmentMaintenanceCommand) Help() string {
	return `Show single equipment.`
}

// Run show equipment
func (c *EquipmentMaintenanceCommand) Run(args []string) int {
	sess, err := auth.Restore()
	if err != nil {
		log.Fatalf("Error restoring session: %s", err)
	}

	id := args[0]
	equipment, err := sess.GetEquipment(context.Background(), id)
	if err != nil {
		log.Fatalf("Error fetching equipment: %s", err)
	}

	history, err := sess.MaintenanceHistory(equipment.ID)
	if err != nil {
		log.Fatalf("Error fetching equipment: %s", err)
	}

	for _, entry := range history {
		print(entry, "")
	}

	return 0
}

// Synopsis returns the equipment list command synopsis
func (c *EquipmentMaintenanceCommand) Synopsis() string {
	return "list equipment maintenance"
}
