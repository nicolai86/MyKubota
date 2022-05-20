package mykubota

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	username = os.Getenv("MYKUBOTA_USERNAME")
	password = os.Getenv("MYKUBOTA_PASSWORD")
	shared   *Session
)

func skipIntegrationWithoutConfiguration(t *testing.T) {
	if username == "" {
		t.Skip("missing MYKUBOTA_USERNAME variable")
	}
	if password == "" {
		t.Skip("missing MYKUBOTA_PASSWORD variable")
	}
}

func hasConfiguration() bool {
	return username != "" && password != ""
}

func TestMain(m *testing.M) {
	if hasConfiguration() {
		session, err := New(context.Background(), username, password)
		if err != nil {
			log.Fatalf("expected login to succeed, but didn't: %v", err)
		}
		shared = session
	}
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	session, err := New(context.Background(), username, password)
	if err != nil {
		t.Fatalf("expected login to succeed, but didn't: %v", err)
	}
	_ = session
}

func TestSession_User(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	user, err := shared.User(context.Background())
	if err != nil {
		t.Fatalf("expected oauth user to succeed, but didn't: %v", err)
	}
	_ = user
}

func TestSession_ListEquipment(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	eqs, err := shared.ListEquipment(context.Background())
	if err != nil {
		t.Fatalf("expected list equipment to succeed, but didn't: %v", err)
	}

	if len(eqs) > 0 {
		eq, err := shared.GetEquipment(context.Background(), eqs[0].ID)
		if err != nil {
			t.Fatalf("expected get equipment to succeed, but didn't: %v", err)
		}
		if diff := cmp.Diff(eqs[0], *eq); diff != "" {
			t.Fatalf("expected same equipment, but got diff\n%s", string(diff))
		}
	}
}

func TestSession_Settings(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	settings, err := shared.Settings(context.Background())
	if err != nil {
		t.Fatalf("expected api settings to succeed, but didn't: %v", err)
	}
	_ = settings
}

func TestSession_Categories(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	categories, err := shared.Categories(context.Background())
	if err != nil {
		t.Fatalf("expected api settings to succeed, but didn't: %v", err)
	}
	_ = categories
}

func TestSession_Models(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	models, err := shared.Models(context.Background())
	if err != nil {
		t.Fatalf("expected api settings to succeed, but didn't: %v", err)
	}
	_ = models
}

func TestSession_SearchMachine(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	model, err := shared.SearchMachine(context.Background(), SearchMachineRequest{
		PartialModel: "kx0",
		Serial:       "39381",
		Locale:       "en-CA",
	})
	if err != nil {
		t.Fatalf("expected api settings to succeed, but didn't: %v", err)
	}
	_ = model
}
