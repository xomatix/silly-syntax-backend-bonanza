package authentication

import (
	"errors"
	"net/http"
)

func ResolveCookiesFromHeader(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		return "", errors.New("bonanza token is missing")
	}

	return token, nil
}

func ResolveCookiesWithUserFromHeader(r *http.Request) (int64, string, error) {
	cookie, err := ResolveCookiesFromHeader(r)
	if err != nil {
		return 0, "", err
	}

	userID, username, err := DecodeToken(cookie)
	if err != nil {
		return 0, "", err
	}
	return userID, username, nil
}
