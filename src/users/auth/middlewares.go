package auth

import (
	"errors"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"pitstop-api/src/utils"
	"strings"
	"time"
)

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := ExtractToken(r)
		if err != nil {
			utils.Prettier(w, "token invalide (bearer manquant)", nil, http.StatusUnauthorized)
			return
		}
		for _, val := range blackListAccessToken {
			if token == val {
				utils.Prettier(w, "token invalide !", nil, http.StatusUnauthorized)
				return
			}
		}

		_, err = ExtractIdFromRequest(r)
		if err != nil {
			utils.Prettier(w, "token invalide !", nil, http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func ExtractIdFromRequest(r *http.Request) (uint, error) {
	headerToken := r.Header.Get("Authorization")
	if !strings.HasPrefix(headerToken, "Bearer ") {
		return 0, errors.New("token invalide (bearer manquant)")
	}
	headerToken = strings.TrimPrefix(headerToken, "Bearer ")
	token, err := jwt.ParseString(headerToken, jwt.WithKeySet(keySet), jwt.UseDefaultKey(true))
	if err != nil {
		return 0, err
	}
	id, ok := token.Get("id")
	if !ok {
		return 0, errors.New("token invalide")
	}

	if !time.Now().UTC().Before(token.Expiration()) {
		return 0, errors.New("token expiré")
	}

	floatId := id.(float64)
	if floatId < 1 {
		return 0, errors.New("token invalide")
	}

	return uint(floatId), nil
}

func ExtractToken(r *http.Request) (string, error) {
	headerToken := r.Header.Get("Authorization")
	if !strings.HasPrefix(headerToken, "Bearer ") {
		return "", errors.New("token invalide (bearer manquant)")
	}
	headerToken = strings.TrimPrefix(headerToken, "Bearer ")

	return headerToken, nil
}

func ExtractUsernameFromRequest(r *http.Request) (string, error) {
	headerToken := r.Header.Get("Authorization")
	if !strings.HasPrefix(headerToken, "Bearer ") {
		return "", errors.New("token invalide (bearer manquant)")
	}
	headerToken = strings.TrimPrefix(headerToken, "Bearer ")
	token, err := jwt.ParseString(headerToken, jwt.WithKeySet(keySet), jwt.UseDefaultKey(true))
	if err != nil {
		return "", err
	}
	username, ok := token.Get("username")
	if !ok {
		return "", errors.New("token invalide")
	}

	if !time.Now().UTC().Before(token.Expiration()) {
		return "", errors.New("token expiré")
	}

	return username.(string), nil
}

func ExtractUsernameFromToken(tokenStr string) (string, error) {

	token, err := jwt.ParseString(tokenStr, jwt.WithKeySet(keySet), jwt.UseDefaultKey(true))
	if err != nil {
		return "", err
	}

	username, ok := token.Get("username")
	if !ok {
		return "", errors.New("token invalide")
	}

	if !time.Now().UTC().Before(token.Expiration()) {
		return "", errors.New("token expiré")
	}

	return username.(string), nil
}
func ExtractIdFromToken(tokenStr string) (uint, error) {

	token, err := jwt.ParseString(tokenStr, jwt.WithKeySet(keySet), jwt.UseDefaultKey(true))
	if err != nil {
		return 0, err
	}

	id, ok := token.Get("id")
	if !ok {
		return 0, errors.New("token invalide")
	}

	if !time.Now().UTC().Before(token.Expiration()) {
		return 0, errors.New("token expiré")
	}

	return uint(id.(float64)), nil
}
