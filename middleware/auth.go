package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	JWT_ISSUER     string = "https://securetoken.google.com/barrenschat-27212"
	JWT_AUD        string = "barrenschat-27212"
	PUBLIC_KEY_URL string = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
)

func AuthUser(param string) (jwt.MapClaims, error) {
	var publicKey *rsa.PublicKey
	var tok *jwt.Token
	var err error
	var publicPEM map[string]string
	var respBody []byte

	resp, err := http.Get(PUBLIC_KEY_URL)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = json.Unmarshal(respBody, &publicPEM)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Loop over the public keys trying to validate
	for _, pem := range publicPEM {
		tok, err = jwt.Parse(param, func(token *jwt.Token) (interface{}, error) {
			publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			return publicKey, err
		})

		// We parsed the token w/o error so move on.
		if err == nil {
			break
		}
	}

	// Check if we had any errors validating the token
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Validate the claims
	if iss, ok := tok.Claims.(jwt.MapClaims)["iss"].(string); ok {
		if iss != JWT_ISSUER {
			return nil, errors.New("Invalid iss claim")
		}
	} else {
		return nil, errors.New("Iss not present in claims")
	}

	if aud, ok := tok.Claims.(jwt.MapClaims)["aud"].(string); ok {
		if aud != JWT_AUD {
			return nil, errors.New("Invalid aud claim")
		}
	} else {
		return nil, errors.New("Aud not present in claims")
	}

	return tok.Claims.(jwt.MapClaims), err
}
