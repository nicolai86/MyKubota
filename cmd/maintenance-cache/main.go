package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nicolai86/mykubota"
	"golang.org/x/sync/semaphore"
)

func main() {
	ctx := context.Background()
	c := mykubota.New("en-CA")

	models, err := c.ListModels(ctx)
	if err != nil {
		log.Fatalf("unable to list models: %v", err)
	}
	log.Printf("fetching maintenance schedule for %d models\n", len(models))

	mu := sync.Mutex{}
	var maintenanceByModel = map[string][]mykubota.Maintenance{}
	sem := semaphore.NewWeighted(20)

	for _, model := range models {
		go func(model mykubota.Model) {
			sem.Acquire(ctx, 1)
			defer sem.Release(1)
			
			schedule, err := c.MaintenanceSchedule(model.Model)
			if err != nil {
				log.Printf("skipping model %q due to error: %v\n", model.Model, err)
			}
			
			mu.Lock()
			maintenanceByModel[model.Model] = schedule
			mu.Unlock()
		}(model)
	}
	sem.Acquire(ctx, 20)

	year, month, day := time.Now().Date()
	payload := bytes.Buffer{}
	if err := json.NewEncoder(&payload).Encode(maintenanceByModel); err != nil {
		log.Fatalf("unable to encode maintenance snapshot: %v", err)

	}
	if err := os.WriteFile(fmt.Sprintf("snapshot-%d%d%d.json", year, month, day), payload.Bytes(), 0644); err != nil {
		log.Fatalf("unable to write snapshot file; %v", err)
	}
}
