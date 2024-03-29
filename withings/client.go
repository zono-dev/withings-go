// license that can be found in the LICENSE file.

// Package withings is UNOFFICIAL sdk of withings API for Go client.
//
package withings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

const (
	authURL          = "http://account.withings.com/oauth2_user/authorize2"
	tokenURL         = "https://wbsapi.withings.net/v2/oauth2"
	defaultTokenFile = ".access_token.json"
)

// Client type
type Client struct {
	Client       *http.Client
	Conf         *oauth2.Config
	Token        *oauth2.Token
	Timeout      time.Duration
	MeasureURL   string
	MeasureURLv2 string
	SleepURLv2   string
}

// ClientOption type for to customize http.Client
type ClientOption func(*http.Client) error

// AuthorizeOffline provides oauth2 authorization for withings in CLI.
// See example/main.go to know the detail.
func AuthorizeOffline(conf *oauth2.Config) (*oauth2.Token, error) {

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	fmt.Printf("URL to authorize:%s\n", url)

	var grantcode string
	fmt.Printf("Open url your browser and Enter your grant code here.\n Grant Code:")
	fmt.Scan(&grantcode)

	token, err := conf.Exchange(newOauthContext(), grantcode)
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
	c.Token = &oauth2.Token{}
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
func (c *Client) ReadToken(path2file string) (*oauth2.Token, error) {
	t, err := readToken(path2file)
	if err != nil {
		return nil, err
	}
	c.Token = t
	c.Client = GetClient(c.Conf, c.Token)
	return c.Token, nil
}

func readToken(path2file string) (*oauth2.Token, error) {
	token := &oauth2.Token{}
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

func saveToken(t *oauth2.Token, path2file string) error {
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
func (c *Client) RefreshToken() (*oauth2.Token, bool, error) {
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

func refreshToken(conf *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := (conf.TokenSource(newOauthContext(), token).Token())
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
func GetNewConf(cid, secret, redirectURL string) oauth2.Config {
	scopes := []string{ScopeActivity, ScopeMetrics, ScopeInfo}
	conf := oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     cid,
		ClientSecret: secret,
		Scopes:       []string{strings.Join([]string(scopes), ",")},
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL,
			TokenURL:  tokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	return conf
}

// GetClient returns *http.Client which based on conf, token.
func GetClient(conf *oauth2.Config, token *oauth2.Token) *http.Client {
	client := oauth2.NewClient(context.Background(), conf.TokenSource(newOauthContext(), token))
	return client
}

// PrintToken print token information.
func (c *Client) PrintToken() {
	printToken(c.Token)
}

func printToken(t *oauth2.Token) {
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

func printConf(conf *oauth2.Config) {
	fmt.Printf("RedirectURL: %v\n", conf.RedirectURL)
	fmt.Printf("ClientID: %v\n", conf.ClientID)
	fmt.Printf("ClientSecret: %v\n", conf.ClientSecret)
	fmt.Printf("Scopes: %v\n", conf.Scopes)
	fmt.Printf("Endpoint(AuthURL): %v\n", conf.Endpoint.AuthURL)
	fmt.Printf("Endpoint(TokenURL): %v\n", conf.Endpoint.TokenURL)
}

// newOauthContext returns context.Context with
// custom http client for withings's access and refresh tokens endpoints.
func newOauthContext() context.Context {
	c := &http.Client{Transport: &oauthTransport{}}
	return context.WithValue(context.Background(), oauth2.HTTPClient, c)
}

// oauthTransport is making custom request and response for withings api.
type oauthTransport struct{}

// RoundTrip customize request and response for withings api.
func (t *oauthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := interceptRequest(r); err != nil {
		return nil, err
	}

	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if err := interceptResponse(res); err != nil {
		return nil, err
	}

	return res, nil
}

// interceptRequest sets action=requesttoken param.
// this param is required for withings api, but not oauth specification.
func interceptRequest(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("cannot parse request form: %v", err)
	}

	req.PostForm.Set("action", "requesttoken")
	encoded := req.PostForm.Encode()
	req.Body = ioutil.NopCloser(strings.NewReader(encoded))
	req.ContentLength = int64(len(encoded))

	return nil
}

// interceptResponse flattens response.
// withings's response body is nested.
// example)
// from:
//     {
//       "status": 0,
//       "body": {
//         "userid": "363",
//         "access_token": "a075f8c14fb8df40b08ebc8508533dc332a6910a",
//         "refresh_token": "f631236f02b991810feb774765b6ae8e6c6839ca",
//         "expires_in": 10800,
//         "scope": "user.info,user.metrics",
//         "csrf_token": "PACnnxwHTaBQOzF7bQqwFUUotIuvtzSM",
//         "token_type": "Bearer"
//       }
//     }
// to:
//     {
//       "userid": "363",
//       "access_token": "a075f8c14fb8df40b08ebc8508533dc332a6910a",
//       "refresh_token": "f631236f02b991810feb774765b6ae8e6c6839ca",
//       "expires_in": 10800,
//       "scope": "user.info,user.metrics",
//       "csrf_token": "PACnnxwHTaBQOzF7bQqwFUUotIuvtzSM",
//       "token_type": "Bearer"
//     }
func interceptResponse(res *http.Response) error {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return fmt.Errorf("cannot read response body: %v", err)
	}

	var withingsRes struct {
		Status int             `json:"status"`
		Body   json.RawMessage `json:"body"`
		Error  string          `json:"error,omitempty"`
	}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&withingsRes)
	if err != nil || withingsRes.Error != "" {
		return &oauth2.RetrieveError{
			Response: res,
			Body:     body,
		}
	}

	res.Body = ioutil.NopCloser(bytes.NewReader(withingsRes.Body))
	res.ContentLength = int64(len(withingsRes.Body))

	return nil
}
