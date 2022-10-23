package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/nicolai86/mykubota"
	"golang.org/x/oauth2"
)

func Persist(s *mykubota.Session) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/.mykubota", home), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	
	return json.NewEncoder(f).Encode(*s.Token)
}

func Restore() (*mykubota.Session, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(fmt.Sprintf("%s/.mykubota", home))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	var token *oauth2.Token
	if err := json.NewDecoder(f).Decode(&token); err != nil {
		return nil, err
	}
	
	return mykubota.New("en-CA").SessionFromToken(context.Background(), token)
}