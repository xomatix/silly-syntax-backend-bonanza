package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xomatix/silly-syntax-backend-bonanza/api/types"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_functions"
	querygenerators "github.com/xomatix/silly-syntax-backend-bonanza/database/queryGenerators"
)

type AuthenticationLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func InitAuthenticationLoginRoutes(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	resp := types.ResponseMessage{
		Success: true,
		Data:    "",
		Message: "success",
	}
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
	}

	var credentials AuthenticationLogin
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	query := querygenerators.SelectQueryCreator{
		CollectionName: "users",
		Filter:         fmt.Sprintf("username = '%s'", credentials.Login),
		Limit:          1,
	}
	generatedSelectQuery, _ := query.GetQuery()
	result, err := database_functions.ExecuteQuery(generatedSelectQuery)
	if err != nil {
		fmt.Println(generatedSelectQuery)
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	if len(result) == 0 {
		resp.Success = false
		resp.Message = "user not found"
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	passwordFromDB := result[0]["password"].(string)
	userIDFromDB := result[0]["id"].(int64)

	if !authentication.CheckPasswordHash(credentials.Password, passwordFromDB) || userIDFromDB == 0 {
		resp.Success = false
		resp.Message = "wrong password"
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	resp.Data, err = authentication.CreateToken(userIDFromDB, credentials.Login)

	if err != nil {
		resp.Success = false
		resp.Message = "error generating token"
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	jsonString, err := json.Marshal(resp)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	cookie := http.Cookie{
		Name:     "bonanza_token",
		Value:    resp.Data.(string),
		MaxAge:   20 * 60,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	w.Write(jsonString)
}
