package cmds

import (
    "context"
    "fmt"
    "github.com/mitchellh/cli"
    "log"
    "reflect"

    "github.com/nicolai86/mykubota"
    "github.com/nicolai86/mykubota/cmd/cli/internal/auth"
)

type EquipmentShowCommand struct{}

var _ cli.Command = &EquipmentShowCommand{}

// Help prints equipment show help
func (c *EquipmentShowCommand) Help() string {
    return `Show single equipment.`
}

// Run show equipment
func (c *EquipmentShowCommand) Run(args []string) int {
    sess, err := auth.Restore()
    if err != nil {
        log.Fatalf("Error restoring session: %s", err)
    }

    id := args[0]
    equipment, err := sess.GetEquipment(context.Background(), id)
    if err != nil {
        log.Fatalf("Error fetching equipment: %s", err)
    }

    print(*equipment, "")

    return 0
}

func print(val interface{}, prefix string) {
    v := reflect.ValueOf(val)
    t := v.Type()
    for i := 0; i < v.NumField(); i++ {
        fmt.Printf("%s%s: ", prefix, t.Field(i).Name)
        switch nv := v.Field(i).Interface().(type) {
        case mykubota.ManualEntries:
            fmt.Println()
            for _, entry := range nv {
                print(entry, fmt.Sprintf("\t%s", prefix))
            }
        case mykubota.VideoEntries:
            fmt.Println()
            for _, entry := range nv {
                print(entry, fmt.Sprintf("\t%s", prefix))
            }

        case mykubota.EquipmentRestartInhibitStatus:
            fmt.Println()
            print(nv, fmt.Sprintf("\t%s", prefix))
        case mykubota.EquipmentLocation:
            fmt.Println()
            print(nv, fmt.Sprintf("\t%s", prefix))
        case mykubota.EquipmentTelematics:
            fmt.Println()
            if !reflect.ValueOf(nv).IsZero() {
                print(nv, fmt.Sprintf("\t%s", prefix))
            }
        default:
            fmt.Printf("%v\n", nv)
        }
    }
}

// Synopsis returns the equipment list command synopsis
func (c *EquipmentShowCommand) Synopsis() string {
    return "show equipment"
}
