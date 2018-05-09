package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	jwt "github.com/dgrijalva/jwt-go"
)

var publicKey *rsa.PublicKey

func MiddlewareChain(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}
func AuthA(param string) (map[string]string, error) {
	var tok *jwt.Token
	var err error

	resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	var publicPEM map[string]string
	err = json.Unmarshal(respBody, &publicPEM)
	if err != nil {
		log.Println(err.Error())
	}
	for _, v := range publicPEM {
		tok, err = jwt.Parse(param, func(token *jwt.Token) (interface{}, error) {
			publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(v))
			return publicKey, err
		})
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Println(err.Error())
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	claims := tok.Claims

	// TODO: validate claims
	log.Println(reflect.TypeOf(claims))

	cl := make(map[string]string)
	for k, v := range claims.(jwt.MapClaims) {
		log.Println(k, v)
		cl[k] = v.(string)
	}

	return cl, err

}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tok *jwt.Token
		var err error

		resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")

		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)

		var publicPEM map[string]string
		err = json.Unmarshal(respBody, &publicPEM)
		if err != nil {
			log.Println(err.Error())
		}
		for _, v := range publicPEM {
			tok, err = jwt.Parse(r.URL.Query().Get("params"), func(token *jwt.Token) (interface{}, error) {
				publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(v))
				return publicKey, err
			})
			if err == nil {
				break
			}
		}

		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		claims := tok.Claims
		// TODO: validate claims
		log.Println(claims)
		next.ServeHTTP(w, r)
	})
}
