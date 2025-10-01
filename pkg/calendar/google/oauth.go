package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

type tokenSavingSource struct {
	oauth2.TokenSource
	filePath string
}

func (s *tokenSavingSource) Token() (*oauth2.Token, error) {
	tok, err := s.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	if tok.RefreshToken == "" {
		if existing, err := tokenFromFile(s.filePath); err == nil && existing.RefreshToken != "" {
			tok.RefreshToken = existing.RefreshToken
		}
	}

	if err := saveToken(s.filePath, tok); err != nil {
		return nil, err
	}

	return tok, nil
}

func getClient(config *oauth2.Config, tokenFilePath string) (*http.Client, error) {
	tok, err := tokenFromFile(tokenFilePath)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}

		if err := saveToken(tokenFilePath, tok); err != nil {
			return nil, err
		}
	}

	tokenSource := config.TokenSource(context.Background(), tok)
	tokenSource = &tokenSavingSource{
		TokenSource: tokenSource,
		filePath:    tokenFilePath,
	}

	return oauth2.NewClient(context.Background(), tokenSource), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil, err
	}

	return tok, nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}

	return tok, nil
}

func saveToken(path string, token *oauth2.Token) error {
	log.Printf("Saving credential file to: %s\n", path)

	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(token); err != nil {
		return err
	}
	return nil
}
