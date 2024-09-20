package authentication

import (
	"net/http"
)

func SetAuthenticationCookies(w http.ResponseWriter, id int64, username string) {

	token, err := CreateToken(id, username)
	if err != nil {
		return
	}

	cookie := http.Cookie{
		Name:     "bonanza_token",
		Value:    token,
		MaxAge:   20 * 60,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}
