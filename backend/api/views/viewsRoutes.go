package views

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"silly-syntax-backend-bonanza/api/types"
	"silly-syntax-backend-bonanza/database"
	"silly-syntax-backend-bonanza/database/authentication"
	"silly-syntax-backend-bonanza/database/permissions"
	querygenerators "silly-syntax-backend-bonanza/database/queryGenerators"
)

type ViewQueryBody struct {
	ViewName string `json:"viewName"`
}

func InitViewsRoutes(w http.ResponseWriter, r *http.Request) {

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

	userID, username, err := authentication.ResolveCookiesWithUserFromHeader(r)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	var viewQuery ViewQueryBody
	err = json.Unmarshal(body, &viewQuery)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	query := querygenerators.SelectQueryCreator{
		CollectionName: "views",
		Filter:         fmt.Sprintf("name = '%s'", viewQuery.ViewName),
		Limit:          1,
	}
	generatedSelectQuery, _ := query.GetQuery()
	result, err := database.ExecuteQuery(generatedSelectQuery)
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
		resp.Message = "view not found"
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	resolvedQuery := result[0]["query"].(string)
	resolvedQuery = permissions.ResolvePermissionsMacro(resolvedQuery, userID)

	viewResult, err := database.ExecuteQuery(resolvedQuery)

	if err != nil {
		resp.Success = false
		resp.Message = "error generating token"
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	resp.Data = viewResult
	jsonString, err := json.Marshal(resp)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		jsonStr, _ := json.Marshal(resp)
		w.Write(jsonStr)
		return
	}

	if userID != 0 {
		authentication.SetAuthenticationCookies(w, userID, username)
	}
	w.Write(jsonString)
}
