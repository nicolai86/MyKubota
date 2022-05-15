package mykubota

import (
	"log"
	"os"
	"testing"
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

	session, err := New(username, password)
	if err != nil {
		t.Fatalf("expected login to succeed, but didn't: %v", err)
	}
	_ = session
}

func TestSession_User(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	session, err := New(username, password)
	if err != nil {
		t.Fatalf("expected login to succeed, but didn't: %v", err)
	}

	user, err := session.User()
	if err != nil {
		t.Fatalf("expected oauth user to succeed, but didn't: %v", err)
	}
	_ = user
}

func TestSession_Equipment(t *testing.T) {
	skipIntegrationWithoutConfiguration(t)

	session, err := New(username, password)
	if err != nil {
		t.Fatalf("expected login to succeed, but didn't: %v", err)
	}

	eq, err := session.Equipment()
	if err != nil {
		t.Fatalf("expected api equipment to succeed, but didn't: %v", err)
	}
	log.Println(eq)
}