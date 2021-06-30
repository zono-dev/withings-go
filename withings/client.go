// license that can be found in the LICENSE file.

// Package withings is UNOFFICIAL sdk of withings API for Go client.
//
package withings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	//"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

const (
	authURL          = "http://account.withings.com/oauth2_user/authorize2"
	tokenURL         = "https://account.withings.com/oauth2/token"
	defaultTokenFile = ".access_token.json"
)

// Client type
type Client struct {
	Client       *http.Client
	Conf         *Config
	Token        *Token
	Timeout      time.Duration
	MeasureURL   string
	MeasureURLv2 string
	SleepURLv2   string
}

// ClientOption type for to customize http.Client
type ClientOption func(*http.Client) error

// AuthorizeOffline provides oauth2 authorization for withings in CLI.
// See example/main.go to know the detail.
func AuthorizeOffline(conf *Config) (*Token, error) {

	url := conf.AuthCodeURL("state", AccessTypeOffline)

	fmt.Printf("URL to authorize:%s\n", url)

	var grantcode string
	fmt.Printf("Open url your browser and Enter your grant code here.\n Grant Code:")
	fmt.Scan(&grantcode)

	token, err := conf.Exchange(context.Background(), grantcode)
	if err != nil {
		fmt.Println("Failed to oauth2 exchange.")
		return nil, err
	}

	return token, nil
}

// ReadSettings read setting file which is yaml file and returns the settings.
func ReadSettings(path2settings string) map[string]string {
	f, err := os.Open(path2settings)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	var m map[string]string

	if err := d.Decode(&m); err != nil {
		log.Fatal(err)
	}

	return m
}

// getNewContext returns new context that has Timeout settings.
func getNewContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// New returns new client.
// cid is client id, secret and redirectURL are parameters that you got them when you setup withings API.
func New(cid, secret, redirectURL string, options ...ClientOption) (*Client, error) {

	conf := GetNewConf(cid, secret, redirectURL)
	c := &Client{}
	c.Conf = &conf
	c.Token = &Token{}
	c.Client = GetClient(c.Conf, c.Token)
	c.Timeout = 5 * time.Second
	c.MeasureURL = defaultMeasureURL
	c.MeasureURLv2 = defaultMeasureURLv2
	c.SleepURLv2 = defaultSleepURLv2

	for _, option := range options {
		err := option(c.Client)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// SetScope sets scope for oauth2 client.
func (c *Client) SetScope(scopes ...string) {
	c.Conf.Scopes = []string{strings.Join([]string(scopes), ",")}
}

// SetTimeout sets timeout setting for http client.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.Timeout = timeout * time.Second
}

// PrintTimeout print timeout setting.
func (c *Client) PrintTimeout() {
	fmt.Printf("Timeout=%v\n", c.Timeout)
}

// ReadToken read from a file and that token is set to client.
func (c *Client) ReadToken(path2file string) (*Token, error) {
	t, err := readToken(path2file)
	if err != nil {
		return nil, err
	}
	c.Token = t
	c.Client = GetClient(c.Conf, c.Token)
	return c.Token, nil
}

func readToken(path2file string) (*Token, error) {
	token := &Token{}
	file, err := os.Open(path2file)
	if err == nil {
		json.NewDecoder(file).Decode(token)
	}
	return token, err
}

// SaveToken save the token in the file.
func (c *Client) SaveToken(path2file string) error {
	var fname string = path2file
	if fname == "" {
		fname = defaultTokenFile
	}
	return saveToken(c.Token, fname)
}

func saveToken(t *Token, path2file string) error {
	file, err := os.Create(path2file)
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(t)
	if err != nil {
		return err
	}
	return nil
}

// RefreshToken get new token if necessary.
func (c *Client) RefreshToken() (*Token, bool, error) {
	newToken, err := refreshToken(c.Conf, c.Token)
	if err != nil {
		return nil, false, err
	}

	var isNewToken bool = false

	if newToken != c.Token {
		c.Token = newToken
		c.Client = GetClient(c.Conf, c.Token)
		isNewToken = true
	}

	return c.Token, isNewToken, nil
}

func refreshToken(conf *Config, token *Token) (*Token, error) {
	newToken, err := (conf.TokenSource(context.Background(), token).Token())
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

// Do is just call `Do` of http.Client.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	ret, err := c.Client.Do(req)
	return ret, err
}

// GetNewConf returns oauth2.Config with client id, secret, and redirectURL
func GetNewConf(cid, secret, redirectURL string) Config {
	scopes := []string{ScopeActivity, ScopeMetrics, ScopeInfo}
	conf := Config{
		RedirectURL:  redirectURL,
		ClientID:     cid,
		ClientSecret: secret,
		Scopes:       []string{strings.Join([]string(scopes), ",")},
		Endpoint: Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
	return conf
}

// GetClient returns *http.Client which based on conf, token.
func GetClient(conf *Config, token *Token) *http.Client {
	client := conf.Client(context.Background(), token)
	return client
}

// PrintToken print token information.
func (c *Client) PrintToken() {
	printToken(c.Token)
}

func printToken(t *Token) {
	layout := "2006-01-02 15:04:05"
	extraKeys := []string{"access_token", "expires_in", "refresh_token", "scope", "token_type", "userid"}

	fmt.Printf("--Token Information--\n")
	fmt.Printf("AccessToken:%s\n", t.AccessToken)
	fmt.Printf("RefreshToken:%s\n", t.RefreshToken)
	fmt.Printf("ExpiryDate:%s\n", t.Expiry.Format(layout))
	fmt.Printf("TokenType:%s\n", t.TokenType)

	// Extra returns interface{}
	for _, k := range extraKeys {
		switch val := t.Extra(k).(type) {
		case string:
			fmt.Printf("%s:%s\n", k, val)
		case float64:
			fmt.Printf("%s:%g\n", k, val)
		default:
			if val != nil {
				fmt.Println(val)
			}
		}
	}
}

// PrintConf print conf information.
func (c *Client) PrintConf() {
	printConf(c.Conf)
}

func printConf(conf *Config) {
	fmt.Printf("RedirectURL: %v\n", conf.RedirectURL)
	fmt.Printf("ClientID: %v\n", conf.ClientID)
	fmt.Printf("ClientSecret: %v\n", conf.ClientSecret)
	fmt.Printf("Scopes: %v\n", conf.Scopes)
	fmt.Printf("Endpoint(AuthURL): %v\n", conf.Endpoint.AuthURL)
	fmt.Printf("Endpoint(TokenURL): %v\n", conf.Endpoint.TokenURL)
}
