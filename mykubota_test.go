package mykubota

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	username = os.Getenv("MYKUBOTA_USERNAME")
	password = os.Getenv("MYKUBOTA_PASSWORD")
)

func skipIntegrationWithoutConfiguration(t *testing.T) {
	if username == "" {
		t.Skip("missing MYKUBOTA_USERNAME variable")
	}
	if password == "" {
		t.Skip("missing MYKUBOTA_PASSWORD variable")
	}
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

	session, err := New(context.Background(), username, password)
	if err != nil {
		t.Fatalf("expected login to succeed, but didn't: %v", err)
	}

	user, err := session.User(context.Background())
	if err != nil {
		t.Fatalf("expected oauth user to succeed, but didn't: %v", err)
	}
	_ = user
}

func TestSession_ListEquipment(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	session, err := New(context.Background(), username, password)
	if err != nil {
		t.Fatalf("expected login to succeed, but didn't: %v", err)
	}

	eqs, err := session.ListEquipment(context.Background())
	if err != nil {
		t.Fatalf("expected list equipment to succeed, but didn't: %v", err)
	}

	if len(eqs) > 0 {
		eq, err := session.GetEquipment(context.Background(), eqs[0].ID)
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

	session, err := New(context.Background(), username, password)
	if err != nil {
		t.Fatalf("expected login to succeed, but didn't: %v", err)
	}

	settings, err := session.Settings(context.Background())
	if err != nil {
		t.Fatalf("expected api settings to succeed, but didn't: %v", err)
	}
	_ = settings
}