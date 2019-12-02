package siwa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat/go-jwx/jwk"
	"golang.org/x/oauth2/jws"
)

const (
	applePublicKeyURL = "https://appleid.apple.com/auth/keys"
)

// TODO: Do we have to get public key like this?
func getPublicKey(keyID string) (interface{}, error) {
	set, err := jwk.Fetch(applePublicKeyURL)
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return nil, err
	}

	var publicKey interface{}
	for _, key := range set.Keys {
		if key.KeyID() == keyID {
			publicKey, err = key.Materialize()
			if err != nil {
				return nil, err
			}
			return publicKey, nil
		}
	}
	return nil, errors.New("public key not found")
}

func getHeader(idToken string) (*jws.Header, error) {
	s := strings.Split(idToken, ".")
	decoded, err := base64.RawURLEncoding.DecodeString(s[0])
	if err != nil {
		return nil, err
	}
	c := &jws.Header{}
	err = json.NewDecoder(bytes.NewBuffer(decoded)).Decode(c)
	return c, err

}

type claims struct {
	Iss   string `json:"iss"`
	Aud   string `json:"aud"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
	Sub   string `json:"sub"`
	CHash string `json:"c_hash"`
	Nonce string `json:"nonce"`
}

func (c *claims) Valid() error {
	// Verify that the iss field contains https://appleid.apple.com
	if c.Iss != appleCom {
		return errors.New("invalid iss")
	}

	// // Verify that the aud field is the developer’s client_id
	// if c.Aud != clientID {
	// 	return errors.New("invalid aud")
	// }

	// Verify that the time is earlier than the exp value of the token
	now := time.Now().Unix()
	if c.Exp < now {
		return errors.New("invalid exp")
	}

	return nil
}

// ValidateIDToken returns verified IDToken
func ValidateIDToken(idToken, code string, checkNonce func(nonce string) error) (*claims, error) {
	if idToken == "" {
		return nil, errors.New("idToken must be specified")
	}

	c := claims{}

	// validation is executed in this function
	token, err := jwt.ParseWithClaims(idToken, &c, func(t *jwt.Token) (interface{}, error) {
		keyID, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("there is no kid")
		}
		key, ok := keyID.(string)
		if !ok {
			return nil, errors.New("keyID cannot be parsed")
		}
		publicKey, err := getPublicKey(key)
		if err != nil {
			return nil, err
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	// If id_token has nonce, it should be checked.
	if checkNonce != nil {
		if err := checkNonce(c.Nonce); err != nil {
			return nil, err
		}
	}

	// Authorization Endpoint から ID Token が code と共に発行される場合は必須
	// http://openid-foundation-japan.github.io/openid-connect-core-1_0.ja.html#HybridIDToken
	if code != "" {
		alg, ok := token.Header["alg"].(string)
		if !ok {
			return nil, errors.New("no alg field on header")
		}
		chash, err := calculateCHash(alg, code)
		if err != nil {
			return nil, err
		}
		if chash != c.CHash {
			return nil, errors.New("difference chash")
		}
	}

	return &c, nil
}
