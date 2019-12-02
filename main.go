package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
	"github.com/kokukuma/siwa-cli/siwa"
)

var (
	clientID    = "com.kokukuma.siwa.test"
	redirectURI = "https://siwa.kokukuma.com/v1/callback2"
	scope       = "email name"
	teamID      = "QZL2SGYQ3S"
	appleCom    = "https://appleid.apple.com"
	appleKey    = "/Users/kanotatsuya/tmp/siwa/AuthKey_BF4R44V675.p8"
)

func main() {
	state := "dsasdflkajdslke"

	config := siwa.Config{
		TeamID:      teamID,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		AppleKey:    appleKey,
		Scope:       scope,
	}
	authURL := config.CreateAuthURL(state)
	fmt.Println(authURL)

	// open browser
	err := exec.Command("open", authURL).Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	// setting callback & get auth code
	ctx := context.Background()
	ar, err := siwa.GetAuthCode(ctx, state)
	if err != nil {
		fmt.Println(err)
		return
	}

	// get id_token from token endpoint
	token, err := siwa.GetAccessToken(ar.Code, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := siwa.ValidateIDToken(token.IDToken, "", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	spew.Dump(c)
}
