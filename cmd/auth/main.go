package main

import (
	"fmt"
	"os"

	"github.com/zono-dev/withings-go/withings"
)

const (
	tokenFile = "access_token.json"
	layout    = "2006-01-02"
	layout2   = "2006-01-02 15:04:05"
)

var (
	client   *(withings.Client)
	settings map[string]string
)

func auth() {
	var err error
	client, err = withings.New(settings["CID"], settings["Secret"], settings["RedirectURL"])

	if err != nil {
		fmt.Println("Failed to create New client")
		fmt.Println(err)
		return
	}
	// default Scope: user.activity,user.metrics,user.info
	fmt.Println(client.Conf.Scopes)

	// When first time you authorize with Withings API, you need to call AuthorizeOffline.
	if _, err := os.Open(tokenFile); err != nil {
		var e error

		client.Token, e = withings.AuthorizeOffline(client.Conf)
		client.Client = withings.GetClient(client.Conf, client.Token)

		if e != nil {
			fmt.Println("Failed to authorize offline.")
		}
	} else {
		// Client can use token file if it has been in tokenFile path.
		_, err = client.ReadToken(tokenFile)

		if err != nil {
			fmt.Println("Failed to read token file.")
			fmt.Println(err)
			return
		}
	}
}

func tokenFuncs() {
	// Show token
	client.PrintToken()

	// Refresh Token if you need
	_, rf, err := client.RefreshToken()
	if err != nil {
		fmt.Println("Failed to RefreshToken")
		fmt.Println(err)
		return
	}
	if rf {
		fmt.Println("You got new token!")
		client.PrintToken()
	}

	// Save Token if you need
	err = client.SaveToken(tokenFile)
	if err != nil {
		fmt.Println("Failed to RefreshToken")
		fmt.Println(err)
		return
	}
}

func main() {
	settings = withings.ReadSettings(".test_settings.yaml")
	fmt.Println(settings)
	auth()
	tokenFuncs()
}
