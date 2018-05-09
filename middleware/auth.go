package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

func AuthUser(param string) (map[string]string, error) {
	var publicKey *rsa.PublicKey
	var tok *jwt.Token
	var err error
	var publicPEM map[string]string
	var respBody []byte
	var claimsMap = make(map[string]string)

	resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	respBody, _ = ioutil.ReadAll(resp.Body)

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

	// TODO extract more clains to validate further along
	claimsMap["user_id"] = tok.Claims.(jwt.MapClaims)["user_id"].(string)
	return claimsMap, err
}

func ValidateClaims(claims map[string]string) error {
	return nil
}

// func MiddlewareChain(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
// 	for _, mw := range middleware {
// 		h = mw(h)
// 	}
// 	return h
// }

// func Auth(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var tok *jwt.Token
// 		var err error

// 		resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")

// 		if err != nil {
// 			log.Println(err)
// 		}
// 		defer resp.Body.Close()

// 		respBody, err := ioutil.ReadAll(resp.Body)

// 		var publicPEM map[string]string
// 		err = json.Unmarshal(respBody, &publicPEM)
// 		if err != nil {
// 			log.Println(err.Error())
// 		}
// 		for _, v := range publicPEM {
// 			tok, err = jwt.Parse(r.URL.Query().Get("params"), func(token *jwt.Token) (interface{}, error) {
// 				publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(v))
// 				return publicKey, err
// 			})
// 			if err == nil {
// 				break
// 			}
// 		}

// 		if err != nil {
// 			log.Println(err.Error())
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		claims := tok.Claims
// 		// TODO: validate claims
// 		log.Println(claims)
// 		next.ServeHTTP(w, r)
// 	})
// }
