package main

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
	"github.com/kokukuma/siwa-cli/siwa"
)

var (
	// kokukum sample
	clientID    = "com.kokukuma.siwa.test"
	redirectURI = "https://siwa.kokukuma.com/localhost"
	scope       = "email name"
	teamID      = "QZL2SGYQ3S"
	appleKey    = "/Users/kanotatsuya/tmp/siwa/AuthKey_BF4R44V675.p8"
	appleKeyID  = "BF4R44V675"
)

func createRandomString(size int) string {
	alpha := "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(buf)
}

func main() {
	state := createRandomString(18)
	nocne := createRandomString(18)

	// TODO: 必要なものが設定されてなかったらエラーを出す.
	// auth / token endpointそれぞれで最低限必要なものは変わる.
	config := siwa.Config{
		TeamID:      teamID,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		AppleKey:    appleKey,
		AppleKeyID:  appleKeyID,
		Scope:       scope,
	}
	authURL := config.CreateAuthURL(state, nocne)
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
	spew.Dump(ar)
	code := ar.Code

	token, err := siwa.GetAccessToken(code, config)
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
