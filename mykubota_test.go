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
		client := New("en-CA")
		session, err := client.Authenticate(context.Background(), username, password)
		if err != nil {
			log.Fatalf("expected login to succeed, but didn't: %v", err)
		}
		shared = session
	}
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	client := New("en-CA")
	session, err := client.Authenticate(context.Background(), username, password)
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

func TestSession_DeleteEquipment(t *testing.T) {
	t.Skip("TODO")
}

func TestSession_AddEquipment(t *testing.T) {
	t.Skip("TODO")
}

func TestClient_Categories(t *testing.T) {
	t.Parallel()

	categories, err := New("en-CA").ListCategories(context.Background())
	if err != nil {
		t.Fatalf("expected api settings to succeed, but didn't: %v", err)
	}
	if len(categories) == 0 {
		t.Fatalf("expected categories to be populated for en-CA")
	}
}

func TestClient_Models(t *testing.T) {
	t.Parallel()

	models, err := New("en-CA").ListModels(context.Background())
	if err != nil {
		t.Fatalf("expected api settings to succeed, but didn't: %v", err)
	}
	if len(models) == 0 {
		t.Fatalf("expected models to be populated for en-CA")
	}
}

func TestClient_SearchMachine(t *testing.T) {
	t.Parallel()

	model, err := New("en-CA").SearchMachine(context.Background(), SearchMachineRequest{
		PartialModel: "kx0",
		Serial:       "1",
	})

	if err != nil {
		t.Fatalf("expected api settings to succeed, but didn't: %v", err)
	}
	
	if expected := "KX057-4"; model.Model != expected {
		t.Fatalf("expected model %q, got %q", expected, model.Model)
	}
}

func TestClient_GetModelTree(t *testing.T) {
	t.Parallel()

	roots, err := New("en-CA").GetModelTree(context.Background())
	if err != nil {
		t.Fatalf("expected api settings to succeed, but didn't: %v", err)
	}

	if len(roots) == 0 {
		t.Fatalf("expected model tree to be populated")
	}

	for _, root := range roots {
		if len(root.SubCategories) == 0 {
			t.Fatalf("expected every root to have subcategories, but %s has none", root.Name)
		}
	}
}
