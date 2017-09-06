/*
	This code uses the thing on the quick start page.
	It has been partially modified.

	Please check the following quick start page for the original code.

	- Go Quickstart  |  Gmail API  |  Google Developers
		- https://developers.google.com/gmail/api/quickstart/go
*/

package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Songmu/prompter"

	gmail "google.golang.org/api/gmail/v1"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var loadedGmail *gmailConfig

type gmailConfig struct {
	Service *gmail.Service
	Context context.Context
	Cancel  context.CancelFunc
	Client  *http.Client
	Config  *oauth2.Config
}

type gmailData struct {
	ID           string
	InternalDate int64
	Snippet      string
	Date         string
	From         string
	To           string
	Subject      string
	Body         string
}

// GetGmailData is retrieve the target mail list
func GetGmailData(last int64) ([]*gmailData, error) {
	srv, err := getGmailService()
	if err != nil {
		return nil, err
	}
	user := "me"

	list, err := srv.Users.Messages.List(user).Do()
	if err != nil {
		return nil, err
	}

	if len(list.Messages) <= 0 {
		return nil, errors.New("Did not have any messages.")
	}

	dataList := make([]*gmailData, 0, len(list.Messages))
	for _, v := range list.Messages {
		msg, err := srv.Users.Messages.Get(user, v.Id).Format("full").Do()
		if err != nil {
			return nil, err
		}

		mailDate := getTimeFromEpoch(msg.InternalDate)
		lastDate := getTimeFromEpoch(last)
		if !mailDate.After(lastDate) {
			// There are no new target mails.
			break
		}

		// set gmail values
		gmailData := &gmailData{}
		gmailData.ID = msg.Id
		gmailData.InternalDate = msg.InternalDate
		gmailData.Snippet = msg.Snippet

		// get header value
		for _, h := range msg.Payload.Headers {
			switch h.Name {
			case "Date":
				gmailData.Date = h.Value
			case "From":
				gmailData.From = h.Value
			case "To":
				gmailData.To = h.Value
			case "Subject":
				gmailData.Subject = h.Value
			}
		}

		// get body value
		raw := msg.Payload.Body.Data
		for _, p := range msg.Payload.Parts {
			// get only text
			if p.MimeType == "text/plain" {
				raw = p.Body.Data
				break
			}
		}
		dec, err := base64.URLEncoding.DecodeString(raw)
		if err != nil {
			return nil, err
		}
		gmailData.Body = string(dec)

		// append list value
		dataList = append(dataList, gmailData)
	}

	if len(dataList) <= 0 {
		return nil, errors.New("There were no applicable messages.")
	}

	return dataList, err
}

// GetSpeakText is get character string to pass to the speak program
func (g *gmailData) GetSpeakText() string {
	fromText := g.From
	rep := regexp.MustCompile(`<.*>`)
	fromText = rep.ReplaceAllString(fromText, "")
	fromText = strings.Replace(fromText, "\"", "", -1)
	fromText = strings.TrimSpace(fromText)
	result := fmt.Sprintf("%s さんからメールが届きました。", fromText)

	subjectText := g.Subject
	subjectText = strings.TrimSpace(subjectText)
	result = result + fmt.Sprintf("%s という件名で、", subjectText)

	sippetText := g.Snippet
	sippetText = html.UnescapeString(sippetText)
	sippetLen := len([]rune(sippetText))
	if sippetLen >= 140 {
		// limit 140 words
		sippetText = string([]rune(sippetText)[:140])
	}
	sippetText = strings.TrimSpace(sippetText)
	result = result + fmt.Sprintf("%s で始まるメールです。", sippetText)

	return result
}

// getGmailService is get client to use gmail api
func getGmailService() (*gmail.Service, error) {
	// return loaded value
	if loadedGmail != nil {
		if loadedGmail.Service != nil {
			return loadedGmail.Service, nil
		}
	}

	// create new service
	loadedGmail = &gmailConfig{}

	// set time out
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	loadedGmail.Context = ctx
	loadedGmail.Cancel = cancel

	// get oauth2 Config
	config, err := getOauth2Config()
	if err != nil {
		return nil, err
	}
	loadedGmail.Config = config

	// get http client
	client, err := getClient(ctx, config)
	if err != nil {
		return nil, err
	}
	loadedGmail.Client = client

	service, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	loadedGmail.Service = service

	return service, err
}

func getOauth2Config() (*oauth2.Config, error) {
	secret, err := getSecret()
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ./.ggm/client_token.json
	config, err := google.ConfigFromJSON(secret, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	return config, err
}

func getSecret() ([]byte, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dataDir, "client_secret.json")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return b, err
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	tok, err := loadToken()
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		err = saveToken(tok)
	}

	return config.Client(ctx, tok), err
}

// getTokenFilePath generates credential file path/filename.
// It returns the generated credential path/filename.
func getTokenFilePath() (string, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dataDir, "client_token.json"), err
}

// loadToken retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func loadToken() (*oauth2.Token, error) {
	path, err := getTokenFilePath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)

	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(token *oauth2.Token) error {
	path, err := getTokenFilePath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// write file
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(token)
	if err != nil {
		return err
	}

	fmt.Printf("\nSaving credential file to: %s\n\n", path)

	return err
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// display URL on console
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n\n%v\n\n", authURL)

	// no echo + 45 characters or more
	code := (&prompter.Prompter{
		Message: "authorization code",
		Regexp:  regexp.MustCompile(`.{45,}`),
		NoEcho:  true,
	}).Prompt()

	tok, err := config.Exchange(oauth2.NoContext, code)

	return tok, err
}

func getGmailProfile() (*gmail.Profile, error) {
	srv, err := getGmailService()
	if err != nil {
		return nil, err
	}
	user := "me"

	prof, err := srv.Users.GetProfile(user).Do()
	if err != nil {
		return nil, err
	}

	return prof, err
}
