package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type tokenSavingSource struct {
	oauth2.TokenSource
	filePath string
}

func (s *tokenSavingSource) Token() (*oauth2.Token, error) {
	tok, err := s.TokenSource.Token()
	if err == nil {
		saveToken(s.filePath, tok)
	}
	return tok, err
}

func getClient(config *oauth2.Config, tokenFilePath string) (*http.Client, error) {
	tok, err := tokenFromFile(tokenFilePath)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}

		saveToken(tokenFilePath, tok)
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

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)
	return nil
}
