package siwa

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	tokenURL = "https://appleid.apple.com/auth/token"
)

// TokenResp represents token endpoint result
type TokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// GetAccessToken connects to Apple and get AT, RT and IDToken in backchannel
func GetAccessToken(code string, config Config) (*TokenResp, error) {

	secret, err := getSecret(config)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("client_id", config.ClientID)
	values.Add("client_secret", secret)
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("redirect_uri", config.RedirectURI)

	req, err := http.NewRequest(
		"POST",
		tokenURL,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(content))
	tokens := TokenResp{}
	err = json.Unmarshal(content, &tokens)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}

func getSecret(config Config) (string, error) {
	iat := time.Now()
	claims := jwt.MapClaims{
		"iss": config.TeamID,
		"aud": appleCom,
		"exp": iat.Add(time.Hour).Unix(),
		"iat": iat.Unix(),
		"sub": config.ClientID,
	}

	alg := jwt.GetSigningMethod("ES256")

	// create a new token
	token := jwt.NewWithClaims(alg, claims)

	token.Header["kid"] = config.AppleKeyID

	privKey, err := loadP8Key(config.AppleKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(privKey)
}

func loadP8Key(path string) (*ecdsa.PrivateKey, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(raw)
	der := block.Bytes

	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}

	switch key := key.(type) {
	case *ecdsa.PrivateKey:
		return key, nil
	default:
		return nil, fmt.Errorf("Found unknown private key type in PKCS#8 wrapping")
	}
}

func calculateCHash(alg string, code string) (string, error) {
	var digest []byte

	switch alg {
	case "ES256", "PS256", "RS256":
		d := sha256.Sum256([]byte(code))
		//left most 256 bits.. 256/8 = 32bytes
		// no need to validate length as sha256.Sum256 returns fixed length
		digest = d[0:32]
	default:
		return "", fmt.Errorf("calculateCHash: %q algorithm not supported", alg)
	}

	left := digest[0 : len(digest)/2]
	return base64.RawURLEncoding.EncodeToString(left), nil
}
